## Сервис СЛОН - Панель управления
Часть системы СЛОН. Позволяет добавлять цели в базу данных бота.

Стек всего проекта:
- cloud.ru
- postgresql
- redis
- rabbitmq
- grafana, prometheus
- nginx
- telegram api
- docker

# Для запуска:
Добавить переменные окружения:
- `SLON_TOKEN` - токен тг бота
- `PSQL_CONN` - connection string для подключения к postgres
- `REDIS_CONN` - connection string для подключения к redis в формате host:port, например localhost:6370
- `REDIS_PASS` - пароль redis

Выполнить:
1. `docker build -t slon_cp .`
2. `docker run -d slon_cp`
