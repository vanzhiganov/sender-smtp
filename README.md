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

### База данных

*Create SMTP record*

```sql
insert into smtp (project_id, server, port, sender_login, sender_password)
values ('3854d9ce-9e27-11ea-9b45-c42c033a81ea', 'smtp.google.com', 587, 'noreply@google.com', '$3cr3t');
```

*Create template*

```sql
insert into `templates` (project_id, template)
values ('3854d9ce-9e27-11ea-9b45-c42c033a81ea', '<!DOCTYPE HTML PULBLIC "-//W3C//DTD HTML 4.01 Transitional//EN">
<html>
<head>
<meta http-equiv="content-type" content="text/html; charset=UTF-8">
</head>
<body>{{ . }}<br>
<div class="moz-signature"><i><br>
<br>
Regards<br>
Vyacheslav Anzhiganov<br>
<i></div>
</body>
</html>');
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