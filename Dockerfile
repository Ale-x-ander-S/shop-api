# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Копируем файлы проекта
COPY . .

# Скачиваем зависимости
RUN go mod download

# Собираем приложение
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd

# Final stage
FROM alpine:latest

WORKDIR /app

# Копируем бинарный файл из builder
COPY --from=builder /app/main .

# Копируем миграции
COPY --from=builder /app/migrations ./migrations

# Создаем пользователя без прав root
RUN adduser -D -g '' appuser
USER appuser

# Запускаем приложение
CMD ["./main"] 