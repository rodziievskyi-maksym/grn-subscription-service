FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG PATH_TO_MAIN=./cmd/go-genesis-case-task/
RUN go build -o go-genesis-case-task-api $PATH_TO_MAIN

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add curl ca-certificates

RUN adduser -D appuser
USER appuser

COPY --from=builder /build/go-genesis-case-task-api .

COPY --from=builder /build/.env .

EXPOSE 8080

CMD ["./go-genesis-case-task-api"]

