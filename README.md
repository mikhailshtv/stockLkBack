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

## Запуск приложения в docker-контейнере

```sh
# Сборка и запуск

docker-compose up --build
```

После запуска API будет доступен по адресу http://localhost:8080, а Swagger UI — по адресу http://localhost:8080/swagger/index.html
