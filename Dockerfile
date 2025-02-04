# Этап сборки: используем официальный образ Golang для компиляции приложения
FROM golang:1.23 AS builder
WORKDIR /app
# Копируем все файлы из текущей директории в контейнер
COPY . .
# Собираем бинарный файл приложения (результат будет называться "service")
RUN go build -o service .

# Этап выполнения: минимальный образ для запуска приложения
FROM alpine:latest
WORKDIR /app
# Копируем скомпилированное приложение из предыдущего этапа
COPY --from=builder /app/service .
# Открываем порт 8080 для доступа к сервису
EXPOSE 8080
# Команда для запуска приложения
CMD ["./service"]
