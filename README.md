# Сервис для создания и просмотра новостей

## Установка и запуск

1. Клонируем репозиторий

2. Выполняем команду

```sh
# Для локального запуска

# go mod tidy && go build -o app.out ./cmd/app
# ./app.out

make run
```

```sh
# Для запуска через docker-compose

# docker-compose up

make service
```