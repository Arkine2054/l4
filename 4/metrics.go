package main

import (
	"runtime"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

var (
	memAlloc = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_alloc_bytes",
		Help: "Currently allocated memory",
	})

	memHeap = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_memory_heap_bytes",
		Help: "Heap memory in use",
	})

	totalAlloc = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "go_memory_total_alloc_bytes",
		Help: "Total bytes allocated since start",
	})

	gcCount = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "go_gc_count_total",
		Help: "Total number of GC runs",
	})

	lastGCTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_last_time_seconds",
		Help: "Last GC time (unix timestamp)",
	})

	gcPause = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "go_gc_pause_last_seconds",
		Help: "Duration of last GC pause",
	})
)

func registerMetrics() {
	prometheus.MustRegister(
		memAlloc,
		memHeap,
		totalAlloc,
		gcCount,
		lastGCTime,
		gcPause,
	)

	go collectRuntimeMetrics()
}

func collectRuntimeMetrics() {
	var prevGC uint32

	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)

		memAlloc.Set(float64(m.Alloc))
		memHeap.Set(float64(m.HeapAlloc))
		totalAlloc.Add(float64(m.TotalAlloc))

		if m.NumGC > prevGC {
			gcCount.Add(float64(m.NumGC - prevGC))
			prevGC = m.NumGC

			if m.LastGC != 0 {
				lastGCTime.Set(float64(m.LastGC) / 1e9)
				gcPause.Set(float64(m.PauseNs[(m.NumGC-1)%256]) / 1e9)
			}
		}

		time.Sleep(5 * time.Second)
	}
}
