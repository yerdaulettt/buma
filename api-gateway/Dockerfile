FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod ./

COPY . .
RUN go build -o api ./cmd/main.go

FROM alpine:3.23.3

COPY --from=builder app/api /api

ENTRYPOINT ["/api"]