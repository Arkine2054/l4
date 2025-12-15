package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
)

// --- структуры для обмена ---
type Request struct {
	Pattern string   `json:"pattern"`
	Lines   []string `json:"lines"`
	Ignore  bool     `json:"ignore"`
}

type Response struct {
	Matches []string `json:"matches"`
}

// --- worker ---
func runWorker(addr string) {
	http.HandleFunc("/grep", func(w http.ResponseWriter, r *http.Request) {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		out := make(chan string)
		for _, line := range req.Lines {
			go func(l string) {
				target := l
				pattern := req.Pattern
				if req.Ignore {
					target = strings.ToLower(l)
					pattern = strings.ToLower(req.Pattern)
				}
				if strings.Contains(target, pattern) {
					out <- l
				} else {
					out <- ""
				}
			}(line)
		}

		var matches []string
		for range req.Lines {
			if m := <-out; m != "" {
				matches = append(matches, m)
			}
		}
		json.NewEncoder(w).Encode(Response{Matches: matches})
	})

	fmt.Printf("Worker listening on %s\n", addr)
	http.ListenAndServe(addr, nil)
}

// --- coordinator ---
func runCoordinator(pattern string, ignore bool, nodes string, quorum int) {
	nodeList := strings.Split(nodes, ",")
	scanner := bufio.NewScanner(os.Stdin)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	// --- распределяем строки по воркерам ---
	chunks := make([][]string, len(nodeList))
	for i, line := range lines {
		chunks[i%len(nodeList)] = append(chunks[i%len(nodeList)], line)
	}

	req := Request{
		Pattern: pattern,
		Ignore:  ignore,
	}

	results := make(chan []string)
	var wg sync.WaitGroup

	for i, n := range nodeList {
		wg.Add(1)
		go func(node string, linesChunk []string) {
			defer wg.Done()
			if len(linesChunk) == 0 {
				results <- []string{}
				return
			}
			req.Lines = linesChunk
			payload, _ := json.Marshal(req)
			resp, err := http.Post("http://"+node+"/grep", "application/json", bytes.NewReader(payload))
			if err != nil {
				results <- []string{}
				return
			}
			defer resp.Body.Close()

			var res Response
			json.NewDecoder(resp.Body).Decode(&res)
			results <- res.Matches
		}(n, chunks[i])
	}

	final := []string{}
	count := 0
	for range nodeList {
		m := <-results
		if len(m) > 0 {
			final = append(final, m...)
			count++
			if count >= quorum {
				break
			}
		}
	}

	// --- выводим уникальные строки ---
	printed := make(map[string]bool)
	for _, l := range final {
		if !printed[l] {
			fmt.Println(l)
			printed[l] = true
		}
	}

	wg.Wait()
}

// --- main CLI ---
func main() {
	workerFlag := flag.Bool("worker", false, "run in worker mode")
	addr := flag.String("addr", ":9001", "worker listen address")
	nodes := flag.String("nodes", "localhost:9001,localhost:9002,localhost:9003", "comma-separated node addresses")
	quorum := flag.Int("quorum", 2, "quorum count")
	ignoreCase := flag.Bool("i", false, "ignore case")
	flag.Parse()

	if *workerFlag {
		runWorker(*addr)
	} else {
		args := flag.Args()
		if len(args) < 1 {
			fmt.Println("Usage: mygrep [options] PATTERN")
			os.Exit(1)
		}
		pattern := args[0]
		runCoordinator(pattern, *ignoreCase, *nodes, *quorum)
	}
}
