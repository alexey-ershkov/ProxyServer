# ProxyServer
Домашняя работа по курсу "Безопасность" 

Ершков Алексей АПО-31

Технопарк МГТУ им. Баумана, 3 семестр

# Перед запуском
- Выполнить скрипт init.sh в папке /certGen
- Создать базу данных и проинициализировать ее с помощью /sql/init.sql
- Заполнить своими данными .env файл

# Запуск 
```
go mod download
go run main.go
```

# Общая информация
На порту PROXY_PORT (1080 по-умолчанию) - запущен прокси сервер
Пример использования:
```
curl -x http://127.0.0.1 https://mail.ru --cacert <Path to repository>/certGen/ca.key
``` 
На порту REQUEST_PORT (80 по-умолчанию) - сервис для повторных запросов
Запрос на `/req?id=<id>` - повторение запроса по id
Любой другой запрос выведет список сохраненных запросов.
####Например:
Запрос
```
curl http://127.0.0.1/
```
Ответ
```
Saved requests
--------------
1) Host: https://mail.ru/test

GET /test HTTP/1.1
Host: mail.ru
Accept: */*
User-Agent: curl/7.64.1

------ RepeatLink: http://127.0.0.1/req?id=1
------ Example: curl -X GET "http://127.0.0.1/req?id=1"

2) Host: https://mail.ru/test

GET /test HTTP/1.1
Host: mail.ru
Accept: */*
User-Agent: curl/7.64.1

------ RepeatLink: http://127.0.0.1/req?id=2
------ Example: curl -X GET "http://127.0.0.1/req?id=2"

3) Host: http://mail.ru/test

GET /test HTTP/1.1
Host: mail.ru
Accept: */*
User-Agent: curl/7.64.1

------ RepeatLink: http://127.0.0.1/req?id=3
------ Example: curl -X GET "http://127.0.0.1/req?id=3"
```
#Поиск уязвимостей XXE
Происходит при повторных запросах. Если в теле запроса есть `<?xml .*?>`, то после него добавляется 
```
<!DOCTYPE foo [
  <!ELEMENT foo ANY>
  <!ENTITY xxe SYSTEM \"file:///etc/passwd\" >]>
<foo>&xxe;</foo>\n
```

В теле ответа ищется строка вида `root:`, в случае если такая есть то в консоли и теле ответа пишется предупреждение от том, что данный запрос небезопасен