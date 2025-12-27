package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime/debug"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	debug.SetGCPercent(100)

	registerMetrics()

	http.Handle("/metrics", promhttp.Handler())

	log.Println("Server started on :8080")
	log.Println("Metrics: http://localhost:8080/metrics")
	log.Println("pprof:   http://localhost:8080/debug/pprof/")

	log.Fatal(http.ListenAndServe(":8080", nil))
}
