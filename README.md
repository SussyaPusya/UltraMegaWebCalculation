# UltraMegaWebCalculation

**UltraMegaWebCalculation** — это веб-сервис для вычисления математических выражений. Он предоставляет API для обработки выражений и возвращает результат вычислений.

## Установка

1. **Клонируйте репозиторий:**

   ```bash
   git clone https://github.com/SussyaPusya/UltraMegaWebCalculation.git
   cd UltraMegaWebCalculation


2.Убедитесь, что у вас установлен Go версии 1.11 или выше. Используйте следующую команду для установки зависимостей:

  ```bash

  go mod tidy
```
------------------------------------------------------------------------------
## ▶️ Запуск
1. Запустите сервер
Выполните следующую команду для запуска сервера:

```bash
go run cmd/main.go
```
По умолчанию сервер будет запущен на localhost:8080.

2. Проверьте работу сервера
Перейдите по адресу http://localhost:8080/ в вашем браузере, чтобы убедиться, что сервер работает.
-----------------------------
## 📡 Использование API
**Конечная точка**
**POST** `/api/v1/calculate`

🔹 **Параметры запроса**:
`expression` (string): Математическое выражение для вычисления.

🔹 **Пример запроса с использованием** `curl`:
```bash
curl --location 'http://localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
🔹 **Пример ответа**:

```json
{
  "result": 6
}
#code:200

```
🔹 **Пример ошибок которые могу возникнуть путём не правильного запроса**:

 ```  "expression": "2+2*2"``` = code: **200**
 
 ```  "expression": "2+b"``` = code: **500 Internal Server Error** 
 
 ```  "expression": "2+"``` = code: **422 Invalid Expression**



--------------------------------------

## 🧪 Тестирование
Запустите тесты для проверки работоспособности проекта:

```bash
go test ./...
```
---------------
## 📁 Структура проекта
```csharp
UltraMegaWebCalculation/
├── cmd/               # Главный файл для запуска приложения
├── internal/          # Логика приложения
├── pkg/               # Пакеты, доступные для повторного использования
├── go.mod             # Список зависимостей
└── README.md          # Описание проекта
```
--------------
## 📜 Лицензия
Проект распространяется под лицензией [MIT](https://github.com/SussyaPusya/UltraMegaWebCalculation/blob/main/LICENSE).

