package app

import (
	"SLON_tg_bot/src/handlers"
	"github.com/go-telegram/bot"
)

type BotApp struct {
	Bot *bot.Bot
}

func NewBot(token, psqlConnString, redisConnString, redisPass string) (*BotApp, error) {
	res, err := NewResources(psqlConnString, redisConnString, redisPass)
	if err != nil {
		return nil, err
	}

	opts := []bot.Option{
		bot.WithDefaultHandler(handlers.DefaultHandler(res.StateManager, res.Repository)),
		bot.WithMessageTextHandler("/start", bot.MatchTypeExact, handlers.InitHandler),
		bot.WithCallbackQueryDataHandler("addUser", bot.MatchTypePrefix, handlers.CallBackHandlerAddUser(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("editUser", bot.MatchTypePrefix, handlers.CallBackHandlerEditUser(res.StateManager)),
		bot.WithCallbackQueryDataHandler("editSched", bot.MatchTypePrefix, handlers.CallBackHandlerEditSched(res.StateManager)),
		bot.WithCallbackQueryDataHandler("addTags", bot.MatchTypePrefix, handlers.CallBackHandlerEditTags(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("delTags", bot.MatchTypePrefix, handlers.CallBackHandlerDelTags(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("editDelete", bot.MatchTypePrefix, handlers.CallBackHandlerEditDelete(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("sched0", bot.MatchTypePrefix, handlers.CallBackHandlerSched(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("sched1", bot.MatchTypePrefix, handlers.CallBackHandlerSched(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("sched2", bot.MatchTypePrefix, handlers.CallBackHandlerSched(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("myUsers", bot.MatchTypePrefix, handlers.CallBackHandlerMyUsers(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("return", bot.MatchTypePrefix, handlers.CallBackHandlerReturn(res.StateManager)),
		bot.WithCallbackQueryDataHandler("tag0", bot.MatchTypePrefix, handlers.CallBackHandlerAssignTag0(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("tag1", bot.MatchTypePrefix, handlers.CallBackHandlerAssignTag1(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("tag2", bot.MatchTypePrefix, handlers.CallBackHandlerAssignTag2(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("tagReady", bot.MatchTypePrefix, handlers.CallBackHandlerTagReady(res.StateManager)),
		bot.WithCallbackQueryDataHandler("dtag0", bot.MatchTypePrefix, handlers.CallBackHandlerDelTag0(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("dtag1", bot.MatchTypePrefix, handlers.CallBackHandlerDelTag1(res.StateManager, res.Repository)),
		bot.WithCallbackQueryDataHandler("dtag2", bot.MatchTypePrefix, handlers.CallBackHandlerDelTag2(res.StateManager, res.Repository)),
	}

	b, err := bot.New(token, opts...)
	if err != nil {
		return nil, err
	}

	return &BotApp{
		Bot: b,
	}, nil
}
