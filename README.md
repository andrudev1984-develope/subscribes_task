# Сервис подписок (тестовый)
## Результат выполнения задания https://nextcloud.effective-mobile.ru/s/ZtcQb9bdZ6RQkyQ?dir=/&openfile=true

## Ключевые элементы
- cmd/migrations/Init.sql - миграция БД сервиса
- openapi/subscribe/openapi.yaml - OpenApi спецификация
- docker-compose.yaml - Docker compose для локального запуска сервиса (запускаются контейнеры с БД и серверной частью сервиса)

## Запуск
- Выполнить команду "docker compose up" в корне проекта
### HTTP ручки будут доступны по адресам вида: http://localhost:8080/subscribes
### Интерфейс Swagger UI доступен по адресу: http://localhost:8080/docs