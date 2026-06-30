Добавьте в код сервера новый эндпоинт `POST /api/shorten`, который будет принимать в теле запроса JSON-объект `{"url":"<some_url>"}` и возвращать в ответ объект `{"result":"<short_url>"}`.

Запрос может иметь такой вид:
```
POST http://localhost:8080/api/shorten HTTP/1.1
Host: localhost:8080
Content-Type: application/json
{
"url": "https://practicum.yandex.ru"
}
```
Ответ может быть таким:
```
HTTP/1.1 201 OK
Content-Type: application/json
Content-Length: 30
{
"result": "http://localhost:8080/EwHXdJfB"
}
```

Не забудьте добавить тесты на новый эндпоинт, как и на предыдущие.

При реализации задействуйте одну из распространённых библиотек:
- `encoding/json`,
- `github.com/mailru/easyjson`,
- `github.com/pquerna/ffjson`,
- `github.com/labstack/echo`.