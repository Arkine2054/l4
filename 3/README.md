
---

# Проект: Calendar Server

## Структура проекта

```
calendar-server/
├── cmd/
│   └── server/
│       └── main.go           # точка входа сервера
├── internal/
│   ├── calendar/
│   │   ├── calendar.go       # бизнес-логика, storage, CRUD
│   │   └── calendar_test.go  # unit-тесты для сервиса
│   ├── http/
│   │   ├── handler.go        # HTTP-хендлеры
│   │   └── router.go         # маршруты
│   └── workers/
│       ├── cleaner.go        # воркер очистки старых событий
│       ├── reminder_pool.go  # воркер напоминаний
│       └── reminder_pool_test.go
└── README.md
```

---

## README.md

````markdown
# Calendar Server

HTTP-сервер календаря с поддержкой:

- CRUD операций с событиями
- Напоминаний через воркеры
- Очистки старых событий
- Логирования запросов

---
````
## Эндпоинты

### Создание события

```http
POST /create_event
Content-Type: application/json

{
    "user_id": 1,
    "date": "2025-12-26",
    "text": "Meeting",
    "remind_at": "2025-12-26T12:00:00Z"
}
```

**Ответ:**

```json
{
    "result": {
        "ID": 1,
        "UserID": 1,
        "Date": "2025-12-26T00:00:00Z",
        "Text": "Meeting",
        "RemindAt": "2025-12-26T12:00:00Z"
    }
}
```

---

### Получить события на день

```http
GET /events_for_day?user_id=1&date=2025-12-26
```

---

### Обновить событие

```http
POST /update_event
Content-Type: application/json

{
    "id": 1,
    "user_id": 1,
    "date": "2025-12-27",
    "text": "Updated Meeting",
    "remind_at": "2025-12-27T12:00:00Z"
}
```

---

### Удалить событие

```http
POST /delete_event
Content-Type: application/json

{
    "id": 1
}
```

---

## Примеры curl

```bash
# Создание события
curl -X POST http://<server_ip>:8080/create_event \
-H "Content-Type: application/json" \
-d '{"user_id":1,"date":"2025-12-26","text":"Meeting","remind_at":"2025-12-26T12:00:00Z"}'

# Получить события на день
curl "http://<server_ip>:8080/events_for_day?user_id=1&date=2025-12-26"

# Обновить событие
curl -X POST http://<server_ip>:8080/update_event \
-H "Content-Type: application/json" \
-d '{"id":1,"user_id":1,"date":"2025-12-27","text":"Updated Meeting"}'

# Удалить событие
curl -X POST http://<server_ip>:8080/delete_event \
-H "Content-Type: application/json" \
-d '{"id":1}'

curl "http://172.20.92.105:8080/events_for_day?user_id=1&date=2025-12-26"

curl "http://172.20.92.105:8080/events_for_week?user_id=1&date=2025-12-26"

curl "http://172.20.92.105:8080/events_for_month?user_id=1&date=2025-12-26"

```

---

## Тестирование

Запустить unit-тесты для сервиса и воркеров:

```bash
go test ./internal/calendar
go test ./internal/workers
```

---

## Особенности

* События хранятся **в памяти**, база данных не используется.
* Воркеры:

    * **Cleaner**: удаляет события старше 30 дней.
    * **ReminderPool**: асинхронно отправляет напоминания.

```
