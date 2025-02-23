# Документация по проекту

## Описание
Проект реализует распределённую систему вычисления арифметических выражений, разбивая их на отдельные операции, которые могут выполняться параллельно.

## Архитектура
- **Оркестратор**: Принимает выражение, разбивает его на задачи и управляет их выполнением.
- **Агент**: Получает задачи от оркестратора, вычисляет результаты и отправляет их обратно.

## Эндпоинты API
### Оркестратор:
- **POST /api/v1/calculate** – Добавление выражения
- **GET /api/v1/expressions** – Получение списка выражений
- **GET /api/v1/expressions/:id** – Получение выражения по ID
- **GET /internal/task** – Получение задачи агентом
- **POST /internal/task** – Отправка результата задачи агентом

## Переменные среды
- `TIME_ADDITION_MS`, `TIME_SUBTRACTION_MS`, `TIME_MULTIPLICATIONS_MS`, `TIME_DIVISIONS_MS` – Время выполнения операций.
- `COMPUTING_POWER` – Количество горутин в агенте.

## Запуск проекта
```sh
git clone <repo_url>
cd project
make build
make run
```

## Примеры использования
```sh
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{"expression": "2+2*2"}'
```

## Тестирование
```sh
go test ./...
```

## Схема работы
![Architecture Diagram](docs/architecture.png)

---

Проект включает тесты, примеры curl, и подробное описание. Для запуска можно использовать Docker Compose или Makefile. Система масштабируема, легко расширяется добавлением новых агентов и предоставляет понятный API.
# projeck3
