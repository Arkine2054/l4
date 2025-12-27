# Go Memory & GC Metrics Server

HTTP-сервер, публикующий Prometheus-метрики памяти и GC,
используя runtime.ReadMemStats, debug.SetGCPercent и pprof.

## Запуск

```bash
go mod tidy
go run .
