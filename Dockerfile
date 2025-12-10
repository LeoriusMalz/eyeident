#FROM ubuntu:latest
#LABEL authors="lev.maltsev"
#
#ENTRYPOINT ["top", "-b"]

## Используем официальный Go образ
#FROM golang:1.25.5
#
## Создаём рабочую директорию
#WORKDIR /sensorProject
#
## Копируем go.mod и go.sum
#COPY go.mod go.sum ./
#
## Загружаем зависимости
#RUN go mod download
#
## Копируем весь проект
#COPY . .
#
## Сборка сервера
#RUN go build -o server ./cmd/server
#
## Открываем порт 5000
#EXPOSE 5000
#EXPOSE 8080
#
## Команда запуска
#CMD ["./server"]


# Сборка Go
FROM golang:1.25.5 AS builder

WORKDIR /app

# Копируем go.mod и go.sum
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

# Образ
FROM alpine:latest

WORKDIR /app

# Копируем бинарник
COPY --from=builder /app/server ./server

# Копируем static и templates
COPY web/static ./web/static
COPY web/templates ./web/templates
COPY internal ./internal

# Открываем порт
EXPOSE 80
EXPOSE 8080

# Запускаем
CMD ["./server", "--port", "8080"]
