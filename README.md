# URL Shortener

URL Shortener — это приложение для сокращения URL, написанное на языке Go с использованием PostgreSQL для хранения данных. Приложение поддерживает базовую авторизацию (Basic Auth) и может быть запущено с помощью Docker Compose.

## Функциональность

- **POST `/url`**: Сокращает URL и возвращает сокращенный alias. JSON тело запроса: `{"url": <url>, "alias": <alias>}`
- **DELETE `/url/<alias>`**: Удаляет сокращенный URL по alias.
- **GET `/<alias>`**: Перенаправляет на оригинальный URL по alias.

## Запуск проекта

Для запуска приложения используйте Docker Compose. Все необходимые сервисы будут запущены одной командой.

1. **Клонируйте репозиторий:**

    ```sh
    git clone https://github.com/Satori27/url-shortener.git
    cd url-shortener
    ```

2. **Запустите приложение с помощью Docker Compose:**

    ```sh
    docker compose --env-file ./.env up
    ```

    Команда запустит контейнеры для приложения и базу данных PostgreSQL.

## Авторизация

Для выполнения POST и DELETE запросов требуется базовая авторизация. Используйте следующие учетные данные для доступа:

- **Пользователь:** `user`
- **Пароль:** `password`

## Примеры запросов

### Сокращение URL (POST)

**Пример запроса на сокращение `https://www.google.com/` в `google`:**

```sh
curl -X POST http://localhost:9346/url \
-u user:password \
-H "Content-Type: application/json" \
-d '{"url": "https://www.google.com/", "alias": "google"}'
```
Если `alias` будет пустым, то серевер сгенерирует случайный

Если такой `alias` уже существует, то сервер вернёт: `{"status":"Error","error":"url already exists"}
`

**Ожидаемый ответ от сервера:**
`{"status":"OK","alias":"google"}`

### Редирект на URL (GET)
Работу редиректа проверить лучше через браузер, в адресной строке бразуера нужно написать:
```
http://localhost:9346/google
```

**Пример запроса на редирект alias'а `google` через `curl`:**

```sh
curl -X GET -I http://localhost:9346/google
```
**Ожидаемый ответ от сервера:**
```
HTTP/1.1 302 Found
Content-Type: text/html; charset=utf-8
Location: https://www.google.com/
Date: Sat, 03 Aug 2024 13:58:23 GMT
Content-Length: 46
```

### Удаление URL (DELETE)

**Пример запроса на удаление alias'а `google`:**

```sh
curl -X DELETE http://localhost:9346/url/google \
-u user:password \
-H "Content-Type: application/json" 
```
**Ожидаемый ответ от сервера:**
`"url deleted"`
(Если такого alias в базе не существует, тогда вернется `"url not exists"`)

## Тесты
Для написания тестов я использовал библиотеку [mockery](https://github.com/vektra/mockery), [gofakeit](https://github.com/brianvoe/gofakeit) , [testing](https://pkg.go.dev/testing).

**Для запуска тестов локально:**

```
go build -o url-shortener cmd/url-shortener/main.go
./url-shortener
go test ./...
```

