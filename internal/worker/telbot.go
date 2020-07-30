package worker

import (
	"telegram-bot/pkg/smap"

	tgbotapi "github.com/Syfaro/telegram-bot-api"
)

func InitTelBotConsumer(users *smap.STMap, bot *tgbotapi.BotAPI) {
	var ucfg tgbotapi.UpdateConfig = tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60
	upChan, err := bot.GetUpdatesChan(ucfg)
	if err != nil {
		panic(err)
	}
	var text string
	for update := range upChan {
		ID := update.Message.Chat.ID

		if ok := users.Get(ID); ok {
			text = "user exist"
		} else {
			users.Set(ID)
			text = "new user registered"
		}

		msg := tgbotapi.NewMessage(ID, text)
		_, _ = bot.Send(msg)
	}
}

func InitTelBotProducer(users *smap.STMap, bot *tgbotapi.BotAPI, ch <-chan string) {
	for msg := range ch {
		for _, chat := range users.GetAllKeys() {
			ms := tgbotapi.NewMessage(chat, msg)
			_, _ = bot.Send(ms)
		}
	}
}
