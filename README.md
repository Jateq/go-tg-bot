# go-tg-bot
Telegram bot using another github module, with goroutines, that returns random image using unsplash api

in .gitignore i have file named constants, that has BotAcces key and UnsplashAcces token

Note that it is not necessary to use mutex when using channels, you don't need both at the same time, nevertheless i needed to show that i can work with goroutines, that's why kept evertything
