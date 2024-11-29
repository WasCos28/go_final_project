FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o todo-app main.go

ENV TODO_PORT=7540 \
    TODO_DBFILE=/app/scheduler.db \
    TODO_PASSWORD=secret

EXPOSE ${TODO_PORT}

CMD ["sh", "-c", "./todo-app -port ${TODO_PORT}"]