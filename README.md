# CalcServer
## Краткое описание
Сервис для подсчёта арифметических выражений
***
## Установка и работа
1. Клонировать репозиторий: git clone https://github.com/Varman56/CalcServer
2. Запустить проект: go run .\cmd\main.go
***
## Использование
### Инструкции
Сервер запускается локально, слушает порт 8080.
На вход принимает POST запрос на адрес /api/v1/calculate, вместе с json в формате:
```
{
    "expression": "выражение, которое ввёл пользователь"
}
```
В ответ также приходит json:
1. В случае успешного вычисления выражения:
```
{
    "result": "результат выражения"
}
```
с кодом 200
2. В случае неудачи:
```
{
    "error": "Сообщение об ошибке"
}
```
с сообщением о причине ошибки, а также, в зависимости от вида ошибки, коды 422/500
***
### Примеры
Опишем несколько рабочих запросов:
1. ```curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"22*2\"}" http://localhost:8080/api/v1/calculate```
/- Получим ответ: ```{"result":44}``` - код 200
2. ```curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"22/0\"}" http://localhost:8080/api/v1/calculate```
/- Получим ответ: ```{"error":"division by zero"}``` - код 422
3. ```curl -X POST -H "Content-Type: application/json" http://localhost:8080/api/v1/calculate```
/- Получим ответ: ```{"error":"invalid json request"}``` - код 500
4. ```curl -X POST -H "Content-Type: application/json" -d "{\"expression\": \"((22.2/2)*3)*(-7)\"}" http://localhost:8080/api/v1/calculate```
/- Получим ответ: ```{"result":-233.09999999999997}``` - код 200 
## Принцип работы
### Калькулятор
Реализовывает интерфейс выражения, в нашем случае - арифметического. В начале строка разбивается на токены (отдельные части выражения), затем рекурсивно (по скобкам) высчитывается, через дополнительную функцию подсчёта операций одного приоритета
### Сервер
Принимает POST-запрос, пытается его обработать. Отлавливает все ошибки, типизирует их и возвращает json-ом с описанием. В случае хорошей работы - отсылает результат выражения, также в json формате
