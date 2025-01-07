package app

import (
	"SLON_tg_bot/src/handlers"
	"github.com/go-telegram/bot"
)

type BotApp struct {
	Bot *bot.Bot
	//Resources *Resources
}

func NewBot(token string) (*BotApp, error) {
	res := NewResources()

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler(res.StateManager, res.Repository)),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.InitHandler),
		bot.WithCallbackQueryDataHandler("addUser", bot.MatchTypePrefix, handlers.CallBackHandlerAddUser(res.StateManager)),
		bot.WithCallbackQueryDataHandler("editUser", bot.MatchTypePrefix, handlers.CallBackHandlerEditUser(res.StateManager)),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	return &BotApp{
		Bot: b,
		//Resources: NewResources(),
	}, nil
}

//func (b *BotApp) RegisterHandlers() {
//	b.Bot.RegisterHandler(bot.handlertype)
//}
