# Time Tracker API

Этот проект реализует тайм-трекер с использованием языка программирования Golang и базы данных PostgreSQL. 

## Функциональность

- **Получение данных пользователей**:
  - Фильтрация по всем полям.
  - Пагинация.
  - Получение трудозатрат по пользователю за период с сортировкой от большей к меньшей затрате времени.
- **Управление задачами**:
  - Начать отсчет времени по задаче для пользователя.
  - Закончить отсчет времени по задаче для пользователя.
- **Управление пользователями**:
  - Удаление пользователя.
  - Изменение данных пользователя.
  - Добавление нового пользователя по номеру паспорта.


## Требования

- Go 1.19+
- Docker
- Docker Compose

## Структура проекта

- **cmd/**: точка входа приложения.
- **internal/delivery/**: слой доставки, реализующий REST API.
- **internal/repository/**: слой доступа к данным, работа с базой данных PostgreSQL.
- **internal/service/**: слой бизнес-логики.
- **internal/models/**: определения моделей данных.

## Настройки сервиса находяться в файле .env

1) `Docker_Port=8000`- Порт, на котором будет работать ваше Go-приложение

2) `DB_USER=postgres`- Имя пользователя для подключения к базе данных PostgreSQL

3) `DB_PASSWORD=yourpassword`- Пароль для подключения к базе данных PostgreSQL

4) `DB_NAME=yourdbname` - Имя базы данных PostgreSQL

5) `DB_PORT=5432`- Порт для подключения к базе данных PostgreSQL

6) `DB_HOST=postgres`- Хост для подключения к базе данных PostgreSQL, указываем имя сервиса из docker-compose

7) `DEBUG=on`- Флаг для включения режима отладки

8) `Service_Url=:8000` - URL, по которому будет доступен сервис

9) `API_URL=http://localhost:8000/info` - URL внешнего API для получения данных о пользователе по номеру паспорта

10) `DB_MIGRATIONS_PATH=file://internal/migrations` - путь, по которому храняться миграции к БД

## API routes

`GET /get_users`: получение данных пользователей с фильтрацией и пагинацией.

`GET /get_user_tasks`: получение трудозатрат по пользователю за период.

`POST /start_user_task`: начать отсчет времени по задаче для пользователя.

`POST /stop_user_task`: закончить отсчет времени по задаче для пользователя.

`POST /add_user`: добавление нового пользователя.

`POST /edit_user`: изменение данных пользователя.

`POST /delete_user`: удаление пользователя.

`GET /swagger/*`: документация Swagger.

## Логирование

Код покрыт debug и info логами для упрощения отладки и мониторинга.

## Запуск web сервиса
Перед запуском необходимо установить в .env файле полный путь до внешнего API ``API_URL``.


1) Запуск в Docker compose:
    -  ``docker-compose up --build``
2) Запуск вручную:
    - ``go run /cmd/main.go``
3) Сборка исполняемого файла:
    - ``go build -o time-tracker ./cmd``



### Helm charts (Автоматическая установка в кластере kubernetes):

1) Установить helm

2) Для установки
```bash
helm install myapp ./myapp --namespace default
```
3) Для удаления
```bash
helm delete myapp
```