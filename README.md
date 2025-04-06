# Shop API

REST API для интернет-магазина на Go с использованием PostgreSQL.


JUST TO TEST

## Требования

- Go 1.21 или выше
- PostgreSQL 15 или выше
- Docker (опционально)

## Установка

1. Клонируйте репозиторий:
```bash
git clone https://github.com/yourusername/shop-api.git
cd shop-api
```

2. Установите зависимости:
```bash
go mod download
```

3. Создайте базу данных PostgreSQL:
```bash
createdb shop
```

4. Примените миграции:
```bash
psql -d shop -f migrations/001_create_products_table.sql
```

## Конфигурация

Создайте файл `.env` в корне проекта со следующими переменными:
```
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=shop
SERVER_PORT=8080
```

## Запуск

```bash
go run cmd/main.go
```

API будет доступно по адресу: http://localhost:8080

## API Endpoints

### Products

- `POST /api/products` - Создать новый продукт
- `GET /api/products` - Получить список всех продуктов
- `GET /api/products/{id}` - Получить продукт по ID
- `PUT /api/products/{id}` - Обновить продукт
- `DELETE /api/products/{id}` - Удалить продукт

## Примеры запросов

### Создание продукта
```bash
curl -X POST http://localhost:8080/api/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Test Product",
    "description": "Test Description",
    "price": 99.99,
    "stock": 100
  }'
```

### Получение всех продуктов
```bash
curl http://localhost:8080/api/products
```

### Получение продукта по ID
```bash
curl http://localhost:8080/api/products/1
```

### Обновление продукта
```bash
curl -X PUT http://localhost:8080/api/products/1 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Product",
    "price": 149.99
  }'
```

### Удаление продукта
```bash
curl -X DELETE http://localhost:8080/api/products/1
```

## Развертывание

### Docker

1. Соберите образ:
```bash
docker build -t shop-api .
```

2. Запустите контейнер:
```bash
docker run -p 8080:8080 --env-file .env shop-api
```

### CI/CD

Проект настроен для автоматического развертывания на VDS Selectel при каждом коммите в ветку main.

## Лицензия

MIT 