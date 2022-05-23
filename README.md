# Algo trade service 

Это концепт сервиса который принимает торговые стратегии на псведо языке, например:
```
any.rsi(14, close, 1440) > 70 and any.rsi(4, close, 5) < 30 and any.rsi(3, close, 1) < 20 and aapl.change_percent(close, 1440) > 2
```
где: 
* any - все символа 
* rsi - индикатор 
* (14, close, 1440) - аргументы индикатора, первый количество свечей, второй тип цены свечи (open/close/high/low), последний тайм фрейм в минутах (1440 - дневка)

Данный алгоритм означает: найти любой символ у которого RSI от 14 свечей, ценой закрытия и тайм фрейме день, значение больше 70 и RSI от 4 свечей, ценой закрытия и тайм фрейме 5 минут, значение меньше 20 и у акции APPLE изменение цены от цены закрытия дневки больше 2 процентов.

## Как работает:
Для символа сперва загружается история на глубину указанную в конфиге, затем подписка на стриминг. После получения новых данных из стриминга для символа идет пересчет на соответствие добавленных стратегий. Если символ удовлетворяет указанным условием то будет сообщение:

 ```
 want to buy ticker: MSFT
 ```

 Так же есть апи из одного роута `/event`, `GET` - посмотреть актуальные стратегии, `POST` - добавить стратегию. Тело для поста:

 ```
{
    "event":"any.rsi(14, close, 1440) > 70 and any.rsi(4, close, 5) < 30 and any.rsi(3, close, 1) < 20 and aapl.change_percent(close, 1440) > 2"
}
 ```

## Как собрать/запустить?

Собрать:
```
go build 
```

Запустить:
записать токен в енв с именем `TINKOFF_TOKEN`
```
./algotrade_service --config=./config.yaml
```
`config.yaml`:
```
provider_tinkoff: - описание для провайдера
  timeout_request: 1m - время таймаута если был превышен лимит запросов
  rate_limit_per_second: 5 - лимит запрос 
  days_with_empty_history: 7 - если нет истории в течение 7 дней, скипать такой символ 
  history_depth: 500 - глубина хранимой истории
  day_offset: - разное поведение на составление глубины на запросы к провайдеру 
    days_for_history_1: 1
    days_for_history_5: 1
    days_for_history_15: 1
    days_for_history_60: 7
    days_for_history_1440: 365
```

## Features

* has no documentation
* has no comment
* has race conditions
* has deadlocks
* has bugs and unexpected behavior 

## PS

Глобально это идея создания многопользовательского сервиса с визуальным языком программирования как [blueprint](https://docs.unrealengine.com/5.0/en-US/blueprints-visual-scripting-in-unreal-engine/#:~:text=The%20Blueprint%20Visual%20Scripting%20system%20in%20Unreal%20Engine%20is%20a,or%20objects%20in%20the%20engine.) 