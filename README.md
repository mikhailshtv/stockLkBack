[![Go Version](https://img.shields.io/github/go-mod/go-version/mikhailshtv/stockLkBack)](https://github.com/mikhailshtv/stockLkBack)
[![Go Report Card](https://goreportcard.com/badge/github.com/mikhailshtv/stockLkBack)](https://goreportcard.com/report/github.com/mikhailshtv/stockLkBack)
[![CI Pipeline](https://github.com/mikhailshtv/stockLkBack/actions/workflows/stockLk-github-ci.yml/badge.svg)](https://github.com/mikhailshtv/stockLkBack/actions/workflows/stockLk-github-ci.yml)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
# stockLkBack (Сервис управления складом)

## Требуется реализовать сервис для управления товарами на складе продавца:

1. Менеджер должен иметь возможность создавать товар в базе, указать наименование товара, количество товара на складе, артикул, закупочную цену, цену продажи. Также он может редактировать все параметры товара и удалять их.
2. Клиент должен иметь возможность создавать заказ, в котором будут перечислены товары и их количество.
3. После создания заказа сервис должен пересчитать количество остатков соответствующих товаров на складе, при отмене заказа снова пересчитать остатки товаров на складе (вернуть удаленные товары на склад).

## Сущности:

1. Товар: 
    Поля: наименование, количество, артикул, закупочная цена, цена продажи.
2. Заказ: 
    Поля: номер, дата создания, дата изменения, статус, сумма. (возможно ещё кем создан, кем изменён).
3. Пользователь
    Роли: клиент, менеджер

## Запуск приложения в docker-контейнере:
docker-compose up --build

