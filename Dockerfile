FROM golang:1.20 AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы go.mod и go.sum и устанавливаем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем остальной исходный код
COPY . .

# Собираем двоичный файл
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main .

# Используем минимальный образ для финального контейнера
FROM scratch

# Копируем собранный бинарник из предыдущего стадии
COPY --from=builder /app/main /main
COPY --from=builder /app/web /web

# Указываем команду запуска
CMD ["/main"]