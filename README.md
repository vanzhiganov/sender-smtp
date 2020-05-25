# sender-smtp
SMTP Sender written in Golang

## Конфигурация

### Переменные окружения

TODO: ...

### Пример файла конфигурации

```yaml
application:
    listen: 0.0.0.0:5555
    secret_key: qwerty
    db: /etc/senderapi/database.db
```

## Использование

### CURL

Запрос для отправки письма в виде HTML

```sh
curl localhost:5555/api/v1 -XPOST \
-H 'X-Secret-Key: qwerty' \
-H 'X-Project-ID: 3854d9ce-9e27-11ea-9b45-c42c033a81ea' \
-d '{"message": "wqeqwewe\ndwdwdw", "subject": "wqewqe", "to": "info@gmail.com", "content-type":"html"}'
```

Запрос для отправки письма в виде текста

```sh
curl localhost:5555/api/v1 -XPOST \
-H 'X-Secret-Key: qwerty' \
-H 'X-Project-ID: 3854d9ce-9e27-11ea-9b45-c42c033a81ea' \
-d '{"message": "wqeqwewe\ndwdwdw", "subject": "wqewqe", "to": "info@gmail.com", "content-type":"plain"}'
```

Ответ

```json
{"id":"6f527846-9dc3-11ea-ae6e-c42c033a81ea"}
```