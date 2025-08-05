[![Go Version](https://img.shields.io/github/go-mod/go-version/mikhailshtv/stockLkBack)](https://github.com/mikhailshtv/stockLkBack)
[![Go Report Card](https://goreportcard.com/badge/github.com/mikhailshtv/stockLkBack)](https://goreportcard.com/report/github.com/mikhailshtv/stockLkBack)
[![CI Pipeline](https://github.com/mikhailshtv/stockLkBack/actions/workflows/stockLk-github-ci.yml/badge.svg)](https://github.com/mikhailshtv/stockLkBack/actions/workflows/stockLk-github-ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# stockLkBack (Сервис управления складом)

**stockLkBack** — это серверное приложение для управления складом и заказами. Сервис позволяет менеджерам и клиентам работать с товарами, создавать и редактировать заказы, а также управлять пользователями с различными ролями.

## Документация API

Полная спецификация и тестирование REST API доступны через Swagger:

- [Swagger UI](http://localhost:8080/swagger/index.html) (после запуска приложения)
- Файлы спецификации: `docs/swagger.yaml`, `docs/swagger.json`

## Основные возможности

1. Менеджер может создавать, редактировать и удалять товары (наименование, количество, артикул, закупочная и продажная цена).
2. Клиент может создавать заказы с перечнем товаров и их количеством.
3. После создания или отмены заказа сервис автоматически пересчитывает остатки товаров на складе.
4. Пользователи имеют роли: клиент, менеджер.

## Сущности

- **Товар**: наименование, количество, артикул, закупочная цена, цена продажи.
- **Заказ**: номер, дата создания, дата изменения, статус, сумма, (опционально: кем создан/изменён).
- **Пользователь**: роли — клиент, менеджер.

## Система логирования и обработки ошибок

### Обзор

В проекте реализована современная система логирования и обработки ошибок, которая обеспечивает:

- **Структурированное логирование** с использованием Zap
- **Централизованную обработку ошибок** с разделением на публичные и внутренние
- **Трейсинг запросов** с уникальными trace ID
- **Автоматическое логирование** всех HTTP запросов
- **Обработку паник** с детальным логированием

### Архитектура

```
Handler Layer (HTTP)
    ↓
Service Layer (Business Logic)
    ↓
Repository Layer (Data Access)
    ↓
Database (PostgreSQL)
```

### Типы ошибок

- `VALIDATION_ERROR` - ошибки валидации (400)
- `NOT_FOUND` - ресурс не найден (404)
- `UNAUTHORIZED` - ошибка авторизации (401)
- `FORBIDDEN` - доступ запрещён (403)
- `CONFLICT` - конфликт данных (409)
- `INTERNAL_ERROR` - внутренняя ошибка сервера (500)
- `DATABASE_ERROR` - ошибка базы данных (500)
- `EXTERNAL_ERROR` - ошибка внешнего сервиса (500)

### Конфигурация

```yaml
logging:
  level: info  # debug, info, warn, error
```

### Примеры использования

```go
// В handlers
if err := ctx.ShouldBindJSON(&userReq); err != nil {
    middleware.HandleError(ctx, errors.NewValidationError("Неверный формат данных", err))
    return
}

// В services
if err != nil {
    logger.GetLogger().Error("failed to create user",
        zap.Error(err),
        zap.String("email", user.Email),
    )
    return nil, errors.NewDatabaseError("ошибка создания пользователя", err)
}
```

### Структура логов

```json
{
  "level": "info",
  "ts": "2024-01-15T10:30:00.000Z",
  "msg": "request started",
  "trace_id": "a1b2c3d4e5f6g7h8",
  "method": "POST",
  "path": "/api/v1/users"
}
```

## Запуск приложения в docker-контейнере

```sh
# Сборка и запуск

docker-compose up --build
```

После запуска API будет доступен по адресу http://localhost:8080, а Swagger UI — по адресу http://localhost:8080/swagger/index.html
