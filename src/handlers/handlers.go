package handlers

import (
	"SLON_tg_bot/src/domains/entities"
	"SLON_tg_bot/src/domains/repositories"
	"SLON_tg_bot/src/state_manager"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func DefaultHandler(
	stateManager state_manager.IStateManager,
	repository repositories.IRepository,
) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.Message.Chat.ID
		state, exists := stateManager.GetState(userId)

		if exists {
			switch state {
			case entities.StateType_WaitingForTargetName:
				// TODO add name check for target
				u := entities.NewDefaultUser(userId, update.Message.Text)
				err := repository.AddUser(u)
				if err != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: userId,
						Text:   "Что-то пошло не так. Попробуйте еще, но лучше не надо",
					})
				} else {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID: userId,
						Text:   "Ok! Юзер записан в нашу базу: " + update.Message.Text,
					})
				}
			default:
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID: userId,
					Text:   "Что-то пошло не так...",
				})
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID: userId,
				Text:   "Нет такой команды, нажмите /start",
			})
		}
		stateManager.ClearState(userId)
	}
}

func InitHandler(ctx context.Context, b *bot.Bot, update *models.Update) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: entities.NewStartupInlineKeyboard(),
	}

	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      update.Message.Chat.ID,
		Text:        "Привет! Выбери одно из доступных действий:",
		ReplyMarkup: kb,
	})
}

func CallBackHandlerAddUser(stateManager state_manager.IStateManager) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userId,
			Text:   "Пришли имя пользователя (то, что начинается на @)",
		})

		stateManager.SetState(userId, entities.StateType_WaitingForTargetName)
	}
}

func CallBackHandlerEditUser(stateManager state_manager.IStateManager) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID: userId,
			Text:   "Пришли имя пользователя (то, что начинается на @)",
		})

		stateManager.SetState(userId, entities.StateType_WaitingForTargetNameToEdit)
	}
}
