package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"

	"gitlab.com/arkine/l4/5/internal/handler"
)

func main() {
	http.HandleFunc("/sum", handler.SumHandler)

	log.Println("listening on :8080")
	log.Println("pprof on :6060")
	go func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	log.Fatal(http.ListenAndServe(":8080", nil))
}
