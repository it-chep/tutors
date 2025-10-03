FROM golang:1.24.5-alpine AS builder

# Устанавливаем зависимости через apk мб --no-cache
RUN apk add \
    ffmpeg \
    make \
    gcc \
    musl-dev \
    nano

WORKDIR /app
COPY . .

# Скачиваем зависимости Go
RUN go mod download

# Экспортируем порты
EXPOSE 8080 7002

# Запускаем приложение
CMD ["go", "run", "cmd/main.go"]