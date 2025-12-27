
---

## Примеры использования

### Запуск worker-ноды

Каждый worker запускается на своём порту:

```bash
./mygrep --worker --addr :9001
```

```bash
./mygrep --worker --addr :9002
```

```bash
./mygrep --worker --addr :9003
```

Worker поднимает HTTP-сервер и принимает запросы на обработку строк.

---

### Использование distributed mygrep

Поиск строк, содержащих `ERROR`, с использованием трёх worker-нод и кворума `2`:

```bash
cat bigfile.txt | ./mygrep ERROR \
  --nodes localhost:9001,localhost:9002,localhost:9003 \
  --quorum 2
```

Параметры:

* `ERROR` — искомый шаблон
* `--nodes` — список worker-нод
* `--quorum` — минимальное количество ответивших нод
* `-i` — игнорирование регистра (опционально)

---

## Сравнительный тест с оригинальной утилитой `grep`

### Тестовое окружение

* ОС: Linux (WSL)
* Go version: `go1.24.1`
* Количество worker-нод: `3`
* Кворум: `2`

---

### Тест с использованием `grep`

```bash
time grep ERROR bigfile.txt > grep.out
```

---

### Тест с использованием distributed `mygrep`

```bash
time cat bigfile.txt | ./mygrep ERROR \
  --nodes localhost:9001,localhost:9002,localhost:9003 \
  --quorum 2 > mygrep.out
```

---

### Проверка корректности результатов

```bash
diff grep.out mygrep.out
```

**Результат:**
Если `diff` не выводит различий — результаты работы `mygrep` полностью совпадают с оригинальной утилитой `grep`.

---