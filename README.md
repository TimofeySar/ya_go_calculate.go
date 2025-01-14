# 🧮 Веб-калькулятор

Проект предоставляет веб-сервис для вычисления математических выражений через HTTP API. Поддерживает базовые арифметические операции и возвращает результаты в формате JSON. 🚀

---

## 📌 Возможности

```markdown
- Поддержка операций: `+`, `-`, `*`, `/`, а также использование скобок.
- Ответы в формате JSON.
- Обработка ошибок и корректные HTTP-коды.
```
## ⚙ Установка и запуск

### 1. Клонирование репозитория

```bash
git clone https://github.com/TimofeySar/ya_go_calculate.go.git
cd ya_go_calculate.go

```
### 2. Запуск приложения

```bash
go run ./cmd/main.go
```

💡 Приложение запустится по адресу: [http://localhost:8080](http://localhost:8080).

## 📚 Использование API

### Эндпоинт: `/api/v1/calculate`

#### Пример успешного запроса:

```bash
curl --location --request POST 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

Успешный ответ:

```json
{
  "result": 6
}
```

#### Пример с некорректным выражением:

```bash
curl --location --request POST 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*a"
}'
```

Ответ с ошибкой:

```json
{
  "error": "некорректный символ: a"
}
```
💡 Код ответа: `422 Unprocessable Entity`.

#### Пример с ошибкой сервера:

Если произойдет внутренняя ошибка, вы получите:

```json
{
  "error": "внутренняя ошибка сервера"
}
```
💡 Код ответа: `500 Internal Server Error`.

## 🌐 Пример использования с Postman

1. Откройте Postman.
2. Создайте новый запрос и выберите метод `POST`.
3. Введите URL: `http://localhost:8080/api/v1/calculate`.
4. Перейдите во вкладку **Body**, выберите `raw` и формат `JSON`.
5. Введите тело запроса:

    ```json
    {
      "expression": "2+2*2"
    }
    ```

6. Нажмите **Send**.
7. Вы получите результат в формате JSON:

    ```json
    {
      "result": 6
    }
    ```

## 🗂️ Структура проекта

```plaintext
calc/
├── main.go          
├── go.mod           
├── calculation/     
│   ├── calc.go      
│   └── calc_test.go 
