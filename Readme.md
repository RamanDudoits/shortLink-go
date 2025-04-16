# ShortLink Service

Сервис для создания и управления короткими ссылками с аутентификацией пользователей и статистикой переходов.

## Особенности

- 🔒 JWT аутентификация
- 📊 Статистика переходов по ссылкам
- 🔗 Автоматическое обнаружение дубликатов
- 🚀 Высокая производительность на Go
- 🐳 Поддержка Docker

## Требования

- Go 1.24.2
- PostgreSQL 13
- (Опционально) Docker

## Установка

### 1. Клонирование репозитория
```bash
git clone https://github.com/RamanDudoits/shortLink-go
cd shortLink-go
```

### 2. Настройка окружения
Создайте файл `.env` в корне проекта:
```env
DB_DSN=postgres://user:password@localhost:5432/shortlink?sslmode=disable
HTTP_ADDR=:8080
JWT_SECRET=your_very_secure_secret
JWT_EXPIRY=24h
```

### 3. Запуск миграций

#### Docker развертывание
```bash
docker-compose up -d
```

Файл `docker-compose.yml` должен содержать:
```yaml

services:
  db:
    image: postgres:13
    restart: always
    volumes:
        - ./tmp/db:/var/lib/postgresql/data
    container_name: shortlink-db
    environment:
      POSTGRES_USER: user
      POSTGRES_PASSWORD: password
      POSTGRES_DB: shortlink
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U user -d shortlink"]
      interval: 5s
      timeout: 5s
      retries: 5

```

```bash
# Установите goose (миграции) и make
go install github.com/pressly/goose/v3/cmd/goose@latest
sudo apt update && sudo apt install -y make

# Примените миграции
export DB_DSN="postgres://user:password@localhost:5432/shortlink?sslmode=disable"
make migrate-up

# Примените откатите миграции
make migrate-down
```

### 4. Запуск сервиса

#### Первый запуск:
```bash
# Установите зависимости
go mod download
```

```bash
go run cmd/api/main.go
```

#### Повторный запуск:
```bash
go run cmd/api/main.go
```

## API Документация

### Аутентификация

#### Регистрация
```http
POST /api/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

#### Вход
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securepassword123"
}
```

Ответ:
```json
{
  "token": "eyJhbGciOi...",
  "user": {
    "id": 1,
    "email": "user@example.com"
  }
}
```

### Работа с ссылками

#### Создание короткой ссылки
```http
POST /api/links
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "url": "https://example.com/very/long/url"
}
```

Ответ:
```json
{
  "id": 1,
  "original_url": "https://example.com/very/long/url",
  "short_code": "abc123",
  "click_count": 0,
  "created_at": "2023-05-20T12:00:00Z"
}
```

#### Получение списка ссылок
```http
GET /api/links
Authorization: Bearer <your_jwt_token>
```
Ответ:
```json
[
    {
        "id": 2,
        "original_url": "https://example.com/very/long/url",
        "short_code": "FWgqV",
        "user_id": 1,
        "click_count": 2,
        "created_at": "2025-04-16T07:49:32.873295Z"
    },
    {
        "id": 1,
        "original_url": "https://google.com",
        "short_code": "4Eli3",
        "user_id": 1,
        "click_count": 0,
        "created_at": "2025-04-16T07:48:46.490322Z"
    }
]
```

#### Обновление ссылки
```http
PUTCH /api/links/{id}/update
Authorization: Bearer <your_jwt_token>
Content-Type: application/json

{
  "link": "https://example.com/very/url"
}
```
Ответ:
```json
{
    "id": 2,
    "original_url": "https://example.com/very/url",
    "short_code": "FWgqV",
    "user_id": 1,
    "click_count": 2,
    "created_at": "2025-04-16T07:49:32.873295Z"
}
```

#### Полученние данных о ссылке
```http
GET /api/links/{id}
Authorization: Bearer <your_jwt_token>
```
Ответ:
```json
{
    "id": 2,
    "original_url": "https://example.com/very/url",
    "short_code": "FWgqV",
    "user_id": 1,
    "click_count": 2,
    "created_at": "2025-04-16T07:49:32.873295Z"
}
```

#### Удаление ссылки
```http
DELETE /api/links/1
Authorization: Bearer <your_jwt_token>
```


#### Редирект
```http
GET /FWgqV
```

## Конфигурация

| Переменная       | По умолчанию                     | Описание                  |
|------------------|----------------------------------|---------------------------|
| DB_DSN           | postgres://...                   | Строка подключения к БД   |
| HTTP_ADDR        | :8080                            | Адрес HTTP сервера        |
| JWT_SECRET       | required                         | Секрет для подписи JWT    |
| JWT_EXPIRY       | 24h                              | Время жизни токена        |

## Разработка

### Генерация миграций
```bash
goose -dir migrations create new_migration_name sql
```