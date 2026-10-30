FROM golang:1.27.3-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o backend ./cmd/main.go

FROM alpine:3.23.3

COPY --from=builder app/backend /backend

ENTRYPOINT ["/backend"]