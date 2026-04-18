FROM golang:1.26-alpine AS builder

RUN apk add --no-cache git ca-certificates && update-ca-certificates

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod download

COPY . .

ARG PATH_TO_MAIN=./cmd/grn-subscription-service/
RUN go build -o grn-subscription-service-api $PATH_TO_MAIN

FROM alpine:latest

WORKDIR /app

RUN apk --no-cache add curl ca-certificates

RUN adduser -D appuser
USER appuser

COPY --from=builder /build/grn-subscription-service-api .

COPY --from=builder /build/.env .

COPY --from=builder /build/web ./web

EXPOSE 8080

CMD ["./grn-subscription-service-api"]
