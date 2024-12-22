# Сервис для вычисления арифметических выражений

Этот проект представляет собой веб-сервис, который позволяет пользователям отправлять арифметические выражения через HTTP и получать результаты их вычислений. Сервис поддерживает базовые арифметические операции: сложение, вычитание, умножение и деление.

# Установка

Для работы с проектом вам понадобится установленный Go (версии 1.16 или выше). 

1. Клонируйте репозиторий:

git clone https://github.com/ваш_логин/ваш_репозиторий.git
cd ваш_репозиторий
text

2. Убедитесь, что все зависимости установлены:

go mod tidy
text

# Запуск сервиса

Чтобы запустить сервис, выполните следующую команду:

go run ./cmd/calc_service/...
text

Сервер будет запущен на порту `8080`.

# Использование

# Эндпоинт

- URL: `/api/v1/calculate`
- Метод: `POST`
- Заголовок: `Content-Type: application/json`
- Тело запроса:
{
"expression": "выражение, которое ввёл пользователь"
}
text

# Примеры запросов

# Успешное выполнение

Запрос для вычисления выражения `2 + 2 * 2`:

curl --location 'http://localhost:8080/api/v1/calculate'
--header 'Content-Type: application/json'
--data '{
"expression": "2 + 2 * 2"
}'
text

Ответ:
{
"result": "6.000000"
}
text

# Ошибка 422 (Неверное выражение)

Запрос с недопустимым символом:

curl --location 'http://localhost:8080/api/v1/calculate'
--header 'Content-Type: application/json'
--data '{
"expression": "2 + a * 2"
}'
text

Ответ:
{
"error": "Expression is not valid"
}
text

# Ошибка 500 (Внутренняя ошибка сервера)

Запрос, который вызывает деление на ноль:

curl --location 'http://localhost:8080/api/v1/calculate'
--header 'Content-Type: application/json'
--data '{
"expression": "2 / 0"
}'
text

Ответ:
{
"error": "division by zero"
}