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


FROM golang:1.25.5 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

WORKDIR /app


COPY --from=builder /app/server ./server

COPY web/static ./web/static
COPY web/templates ./web/templates
COPY internal ./internal
#COPY data ./data

EXPOSE 80
EXPOSE 8080

CMD ["./server", "--port", "8080"]
