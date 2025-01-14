FROM golang:1.22

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем файлы для управления зависимостями и загружаем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь проект
COPY . .

# Сборка бинарного файла
RUN go build -o main ./cmd/main.go

# Указываем порт
EXPOSE 8080

# Команда запуска
CMD ["/main"]