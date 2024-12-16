# ya_go_calculate.go

Это проект калькулятора, который я делал ещё в прошлом спринте, но теперь добавил к нему сервер. В итоге получилось что-то, что умеет принимать запросы на подсчёт математических выражений. Работает через API, всё очень просто.

Как это работает
Сервер принимает запросы на localhost/api/v1/calculate с JSON, в котором есть поле expression. Вы туда пишете выражение, а он пытается его посчитать. Иногда считает, иногда ругается, ну, как бывает.

чуть чуть примеров
ввод
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
вывод
{
  "result": 6
}
ввод
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+*"
}'
вывод
{'error': 'некорректное выражение: недостаточно цифр'}
ввод
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "1/0"
}'
вывод
{'error': 'деление на ноль'}


