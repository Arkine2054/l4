
---

# README.md

## Проект: Простой HTTP API `/sum` с профилировкой и бенчмарками

### Структура проекта

```
D:.
|   go.mod
|   go.sum
|   README.md
|
+---bench
|       sum_bench_test.go
|
+---cmd
|   \---api
|           main.go
|
\---internal
    \---handler
            sum.go
```

* `cmd/api/main.go` — точка входа для HTTP-сервиса.
* `internal/handler/sum.go` — реализация хэндлера `SumHandler`.
* `bench/sum_bench_test.go` — бенчмарки для `SumHandler`.
---

### Запуск API

```powershell
cd cmd/api
go run main.go
```

API слушает на `http://localhost:8080/sum` и принимает JSON вида:

```json
{
    "A": 5,
    "B": 7
}
```

Возвращает:

```json
{
    "Sum": 12
}
```

---

### Профилировка

1.CPU-профиль:

```powershell
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30
```

2. Heap-профиль:

```powershell
go tool pprof http://localhost:6060/debug/pprof/heap
```

**Результаты до исправления `log.Printf`:**

* Время выполнения `SumHandler` занимало ~45% CPU из-за блокирующего `log.Printf`.

---

### Оптимизация

* `log.Printf` заменён на `http.Error` для ошибок:

* Результат: **уменьшение времени на одну итерацию бенчмарка на ~1.8%**, без изменений по памяти.

---

### Бенчмарки
**Сравнение до/после:**

| Метрика   | Старое | Новое | Разница |
| --------- | ------ | ----- | ------- |
| ns/op     | 2783   | 2732  | -1.8%   |
| B/op      | 7112   | 7112  | 0%      |
| allocs/op | 27     | 27    | 0%      |

> Улучшение: `SumHandler` стал чуть быстрее, убраны ненужные блокировки на логирование.

---

### Вывод

* Профилирование CPU/Heap помогло найти узкое место — `log.Printf`.
* Оптимизация хэндлера снизила время обработки запросов.
