FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go build -o health-service main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/health-service .
EXPOSE 8080
CMD ["./health-service"]