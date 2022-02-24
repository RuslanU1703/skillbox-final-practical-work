### Skillbox: Профессия Go-разработчик
## Итоговый проект: Cетевой многопоточный сервис для Statuspage
## Автор: Руслан Юсупов @ruslanU1703

Каталоги **diplom** и **localView** включают в себя API **data** и сервер **myapp**, который получает данные с API, обрабытвает их, и также отдает по endpoint

**data** enpoints:
* /mms - данные по системе MMS
* /static/files/sms.data - данные по системе SMS
* /static/files/voice.data - данные по системе Voice Call
* /static/files/email.data - данные по системе Email
* /static/files/billing.data - данные по системе Billing
* /support - данные по системе поддержки
* /accendent - данные по системе инцидентов
* /test - заглушка для первичной демонстрации StatusPage синтетическими данными
На момент защиты итоговой работы **data** запушен на heroku и доступен по ссылке **https://salty-inlet-33171.herokuapp.com**

**myapp** enpoints:
* /api - обработанные данные с API
На момент защиты итоговой работы **myapp** запушен на heroku и доступен по ссылке **https://glacial-island-82428.herokuapp.com**

Каталог **localView** полностью идентичен **diplom**, за исключением того, что **localView** использует Docker.
Следовательно можно локально посмотреть на работоспособность через **docker-compose.yml** в корне каталога, выполнив
```
docker compose up
```
Подробнее о обработке данных с API можно узнать в ТЗ по дипломному проекту.pdf



