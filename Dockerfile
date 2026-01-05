# ===== STAGE 1: build =====
FROM golang:1.24-bookworm AS builder

WORKDIR /app

RUN apt-get update && apt-get install -y gcc libc6-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o todo-app ./cmd/server

# ===== STAGE 2: runtime =====
FROM debian:bookworm-slim

WORKDIR /app

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /app/todo-app /app/todo-app
COPY web /app/web

EXPOSE 7540

ENV TODO_PORT=7540
ENV TODO_DBFILE=/data/todo.db
ENV TODO_PASSWORD=

CMD ["./todo-app"]
