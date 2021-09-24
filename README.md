##Websocket API Gateway 
Биржа Bitmex предоставляет собственный API интерфейс (REST + Websocket), позволяющий подписываться на различные уведомления об изменениях цены котировок финансовых инструментов


Необходимо реализовать с помощью Gin + Gorilla websocket каркаса шлюз для подписки на сигналы тестового счета биржи Bitmex, при подключению к которому возможно получать оповещения об изменениях в таблице финансовых инструментов (котировках) тестового окружения биржи Bitmex.
Веб-сокет API сервер должен давать возможность клиенту подписаться/отписаться к  каналу, передающего данные о котировках.


Формат сообщения о подписки:
```
JSON: {“action”: “subscribe”, 
       “symbols”: <[]string>}
``` 
, где symbols (опциональный) -  представляет из себя список названий символов торговых инструментов биржи, в случае если поле символа не указано, подписываться к уведомлениям всех символов

Формат сообщения об окончании подписки:

``` 
JSON: {“action”: “unsubscribe”}
``` 

В сообщении веб-сокет канала необходимо передавать следующую информацию



``` 
{
   timestamp: <timestamp>,
   symbol: <symbol_name>,
   price: <lastPrice>
}
``` 
В качестве источника данных использовать подписку “instrument” тестового вебсокет API сервиса Bitmex.
При этом для соединения с веб-сокет сервером Bitmex необходимо использовать только одно подключение, независимо от количества клиентских подписок

##Usage
Start:

```buildoutcfg
$ go build -o ./bin/gateway ./cmd/gateway/main.go

$ ./bin/gateway
```

Listen:

```buildoutcfg

$ gws client -url="ws://0.0.0.0:9090/" -verbose

> {"action":"subscribe"}
```