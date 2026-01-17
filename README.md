# Запуск проекта
docker-compose up --build

Сервер будет доступен по адресу: http://localhost:8080

Миграции
docker-compose exec app goose -dir migrations postgres \
"host=db user=user password=password dbname=chatdb sslmode=disable" up
