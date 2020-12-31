[![Build Status](https://travis-ci.com/Maksimall89/dozor18_bot.svg?branch=master)](https://travis-ci.com/Maksimall89/dozor18_bot) [![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=Maksimall89_dozor18_bot&metric=alert_status)](https://sonarcloud.io/dashboard?id=Maksimall89_dozor18_bot) 
# @dozor18_bot
Telegram бот для проверки кодов приквелов и генерации кодов [@dozor18_bot](https://t.me/dozor18_bot).
Весь список команд бота доступен по команде `/start` или `/help`.

## Настройка
Для работы бота необходимо определить следующие перемерные среды:
* `TelegramBotToken` - token полученный от [@BotFather](https://t.me/BotFather);
* `OwnName` - имя владельца, котрый сможет добавлять коды в бот;
* `ListenPath` - URL по которому бот будет ожидать Webhook;
* `PORT` - порт по которому бот будет ожидать Webhook;
* `DriverNameDB` - имя драйвера для БД (например `postgres`);
* `DATABASE_URL` - connection string для БД;

Далее необходимо настроить Webhook и проверить, что он работает.
* `https://api.telegram.org/bot<TelegramBotToken>/setWebhook?url=https://myapp.com/read<ListenPath>` - установка Webhook;
* `https://api.telegram.org/bot<TelegramBotToken>/getWebhookInfo` - проверка, что webhook работает;
* `https://api.telegram.org/bot<TelegramBotToken>/setWebhook` - удаление webhook в случае ошибки.


