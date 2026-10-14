# Food Delivery — Admin API

REST API администраторской части сервиса доставки еды.  
Реализовано на Go (Chi) согласно принципам Clean Architecture.

## Структура проекта

```
food-delivery-admin/
├── cmd/server/main.go              # Точка входа, DI-конфигурация
├── internal/
│   ├── domain/entities.go          # Доменные сущности (Order, User, Refund)
│   ├── ports/ports.go              # Интерфейсы (репозитории, шлюзы, сервисы)
│   ├── usecase/
│   │   ├── order_usecase.go        # Сценарии: заказы + возврат
│   │   └── user_usecase.go         # Сценарии: пользователи
│   └── infrastructure/
│       ├── repository/
│       │   ├── order_repo.go       # In-memory репозиторий заказов
│       │   └── user_repo.go        # In-memory репозиторий пользователей
│       ├── payment/gateway.go      # Имитация платёжного шлюза
│       └── notification/service.go # Логирующий сервис уведомлений
├── api/handler/
│   ├── order_handler.go            # HTTP-обработчики заказов
│   ├── user_handler.go             # HTTP-обработчики пользователей
│   └── helpers.go                  # writeJSON / writeError
└── docs/openapi.yaml               # OpenAPI 3.0 спецификация
```

## Запуск

```bash
go mod tidy          # загрузка зависимостей
go run ./cmd/server  # запуск сервера на :8080
```

## Endpoints

### Сервис управления заказами
| Метод  | Путь                               | Описание                       |
|--------|------------------------------------|--------------------------------|
| GET    | /admin/orders                      | Список заказов (фильтр ?status)|
| GET    | /admin/orders/{order_id}           | Детали заказа                  |
| PATCH  | /admin/orders/{order_id}/status    | Изменение статуса              |
| POST   | /admin/orders/{order_id}/cancel    | Отмена заказа                  |
| PATCH  | /admin/orders/{order_id}/courier   | Переназначение курьера         |

### Сервис обработки платежей
| Метод  | Путь                               | Описание                       |
|--------|------------------------------------|--------------------------------|
| POST   | /admin/orders/{order_id}/refund    | Возврат средств                |

### Сервис управления пользователями
| Метод  | Путь                               | Описание                       |
|--------|------------------------------------|--------------------------------|
| GET    | /admin/users                       | Список пользователей           |
| GET    | /admin/users/{user_id}             | Пользователь по ID             |
| PATCH  | /admin/users/{user_id}/block       | Блокировка                     |
| PATCH  | /admin/users/{user_id}/unblock     | Разблокировка                  |

## Тестовые данные (seed)

### Заказы
- `order-001` — статус `assigned`, оплачен, курьер `courier-001`
- `order-002` — статус `pending`, оплачен, без курьера
- `order-003` — статус `delivering`, оплачен, курьер `courier-002`

### Пользователи
- `user-001` — Анна Иванова, не заблокирован
- `user-002` — Пётр Сидоров, не заблокирован
- `user-003` — Мария Кузнецова, **заблокирован**
- `courier-001`, `courier-002` — курьеры

## Примеры запросов (Postman / Insomnia)

```bash
# Получить все заказы
GET http://localhost:8080/admin/orders

# Фильтр по статусу
GET http://localhost:8080/admin/orders?status=pending

# Изменить статус
PATCH http://localhost:8080/admin/orders/order-001/status
Body: {"status": "delivering"}

# Отмена заказа (409 — если уже отменён/доставлен)
POST http://localhost:8080/admin/orders/order-002/cancel

# Возврат средств
POST http://localhost:8080/admin/orders/order-001/refund

# Переназначить курьера
PATCH http://localhost:8080/admin/orders/order-002/courier
Body: {"courier_id": "courier-001"}

# Заблокировать пользователя (409 — если уже заблокирован)
PATCH http://localhost:8080/admin/users/user-001/block

# Разблокировать
PATCH http://localhost:8080/admin/users/user-003/unblock
```
