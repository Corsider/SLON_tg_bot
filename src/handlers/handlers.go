package handlers

import (
	"SLON_tg_bot/src/domains/entities"
	"SLON_tg_bot/src/domains/repositories"
	"SLON_tg_bot/src/state_manager"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"regexp"
	"strconv"
)

var tgNameCheck = regexp.MustCompile(`^@[a-zA-Z0-9_]{2,}$`)

func DefaultHandler(
	stateManager state_manager.IStateManager,
	repository repositories.IRepository,
) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		if update == nil {
			return
		}
		userId := update.Message.Chat.ID
		state, exists := stateManager.GetState(userId)

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.NewStartupInlineKeyboard(),
		}

		if exists {
			switch state {
			case entities.StateType_WaitingForTargetName:
				u := entities.NewDefaultUser(userId, update.Message.Text)
				if !nameCheck(update.Message.Text) {
					invalidTgNameMessage(ctx, b, userId)
					break
				}

				existedUser, err := repository.GetSingleByCreatorAndTarget(userId, update.Message.Text)
				if err != nil {
					errorMessage(ctx, b, userId)
					break
				}

				if existedUser != nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      userId,
						Text:        "Такой юзер уже добавлен",
						ReplyMarkup: kb,
					})
					return
				}

				err = repository.AddUser(u)
				if err != nil {
					errorMessage(ctx, b, userId)
				} else {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      userId,
						Text:        "Ok! Юзер записан в нашу базу: " + update.Message.Text,
						ReplyMarkup: kb,
					})
				}
			case entities.StateType_WaitingForTargetNameToEdit:
				if !nameCheck(update.Message.Text) {
					invalidTgNameMessage(ctx, b, userId)
					break
				}
				u, err := repository.GetSingleByCreatorAndTarget(userId, update.Message.Text)
				if err != nil || u == nil {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      userId,
						Text:        "Юзер не найден.",
						ReplyMarkup: kb,
					})
					break
				}

				kb1 := &models.InlineKeyboardMarkup{
					InlineKeyboard: entities.EditInlineKeyboard(),
				}

				// todo add ToFlatUser method
				flatUser := "Юзер: " + u.Target + "\nТип расписания: "
				switch u.Schedule {
				case entities.ScheduleType_ByMessage:
					flatUser += "Триггер на каждое сообщение"
				case entities.ScheduleType_EveryHour:
					flatUser += "Каждый час"
				case entities.ScheduleType_EveryDay:
					flatUser += "Каждый день в 12:00"
				}
				tags := "*нет тегов*"
				if u.Tags != nil {
					tags = *u.Tags
				}
				flatUser += "\nТеги: " + tags

				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      userId,
					Text:        "Выберите, что хотите отредактировать у юзера:\n\n" + flatUser,
					ReplyMarkup: kb1,
				})
				stateManager.SetUser(userId, update.Message.Text)
			case entities.StateType_WaitingForTags:
				target, exists := stateManager.GetUser(userId)
				if !exists {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      userId,
						Text:        "Такого юзера нет.",
						ReplyMarkup: kb,
					})
					break
				}
				// todo check tags spelling
				if len(update.Message.Text) > 500 {
					b.SendMessage(ctx, &bot.SendMessageParams{
						ChatID:      userId,
						Text:        "Слишком большая длина текста тегов",
						ReplyMarkup: kb,
					})
				}
				err := repository.UpdateUserTags(userId, target, update.Message.Text)
				if err != nil {
					errorMessage(ctx, b, userId)
					break
				}
				b.SendMessage(ctx, &bot.SendMessageParams{
					ChatID:      userId,
					Text:        "Теги юзера " + target + " обновлены.",
					ReplyMarkup: kb,
				})
			default:
				errorMessage(ctx, b, userId)
			}
		} else {
			b.SendMessage(ctx, &bot.SendMessageParams{
				ChatID:      userId,
				Text:        "Нет такой команды, нажмите /start или выберите одно из действий ниже:",
				ReplyMarkup: kb,
			})
		}
		stateManager.ClearState(userId)
	}
}

func errorMessage(ctx context.Context, b *bot.Bot, chatId int64) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: entities.NewStartupInlineKeyboard(),
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        "Что-то пошло не так.",
		ReplyMarkup: kb,
	})
}

func invalidTgNameMessage(ctx context.Context, b *bot.Bot, chatId int64) {
	kb := &models.InlineKeyboardMarkup{
		InlineKeyboard: entities.NewStartupInlineKeyboard(),
	}
	b.SendMessage(ctx, &bot.SendMessageParams{
		ChatID:      chatId,
		Text:        "Введено невалидное имя аккаунта телеграмм",
		ReplyMarkup: kb,
	})
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

func CallBackHandlerAddUser(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		users, err := repo.GetUsersByCreator(userId)
		if err != nil {
			errorMessage(ctx, b, userId)
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.NewStartupInlineKeyboard(),
		}

		if len(users) > 4 {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      userId,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        "Нельзя добавить больше 5 юзеров!",
				ReplyMarkup: kb,
			})
			return
		}

		kbRet := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.ReturnInlineKeyboard(),
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userId,
			Text:        "Пришли имя пользователя, которого хочешь добавить (то, что начинается на @)",
			ReplyMarkup: kbRet,
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

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.ReturnInlineKeyboard(),
		}

		b.SendMessage(ctx, &bot.SendMessageParams{
			ChatID:      userId,
			Text:        "Пришли имя пользователя для редактирования (то, что начинается на @)",
			ReplyMarkup: kb,
		})

		stateManager.SetState(userId, entities.StateType_WaitingForTargetNameToEdit)
	}
}

func CallBackHandlerEditSched(stateManager state_manager.IStateManager) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.SchedulesInlineKeyboard(),
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		// todo add current schedule to Text
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      userId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Выбери один из типов расписания:",
			ReplyMarkup: kb,
		})

		stateManager.ClearState(userId)
	}
}

func CallBackHandlerEditTags(stateManager state_manager.IStateManager) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.ReturnInlineKeyboard(),
		}

		// todo add current tags to Text
		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      userId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Пришли теги для уточнения поведения модели через запятую",
			ReplyMarkup: kb,
		})

		stateManager.SetState(userId, entities.StateType_WaitingForTags)
	}
}

func CallBackHandlerEditDelete(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.NewStartupInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			stateManager.ClearState(userId)
			return
		}

		err := repo.RemoveUser(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			stateManager.ClearState(userId)
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      userId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Юзер " + target + " удален.",
			ReplyMarkup: kb,
		})

		stateManager.ClearState(userId)
	}
}

func CallBackHandlerSched(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			stateManager.ClearState(userId)
			return
		}

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.NewStartupInlineKeyboard(),
		}

		switch update.CallbackQuery.Data {
		case "sched0":
			err := repo.UpdateUserSched(userId, target, entities.ScheduleType_ByMessage)
			if err != nil {
				errorMessage(ctx, b, userId)
				stateManager.ClearState(userId)
				return
			}
		case "sched1":
			err := repo.UpdateUserSched(userId, target, entities.ScheduleType_EveryHour)
			if err != nil {
				errorMessage(ctx, b, userId)
				stateManager.ClearState(userId)
				return
			}
		case "sched2":
			err := repo.UpdateUserSched(userId, target, entities.ScheduleType_EveryDay)
			if err != nil {
				errorMessage(ctx, b, userId)
				stateManager.ClearState(userId)
				return
			}
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Расписание юзера " + target + " обновлено на " + func(s string) string {
				switch s {
				case "sched0":
					return "Триггер на каждое сообщение"
				case "sched1":
					return "Каждый час"
				case "sched2":
					return "Каждый день в 12:00"
				}
				return ""
			}(update.CallbackQuery.Data),
			ReplyMarkup: kb,
		})

		stateManager.ClearState(userId)
	}
}

func CallBackHandlerMyUsers(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.NewStartupInlineKeyboard(),
		}

		myTargets, err := repo.GetUsersByCreator(userId)
		if err != nil {
			errorMessage(ctx, b, userId)
			return
		}

		if myTargets == nil || len(myTargets) == 0 {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:      userId,
				MessageID:   update.CallbackQuery.Message.Message.ID,
				Text:        "У вас пока нет добавленных юзеров!",
				ReplyMarkup: kb,
			})
			return
		}

		resultMsg := "Твои юзеры (" + strconv.Itoa(len(myTargets)) + "/5)" + ":\n"
		for _, t := range myTargets {
			resultMsg += "\n===\nЮзер: " + t.Target + ",\nТип расписания: "
			switch t.Schedule {
			case entities.ScheduleType_ByMessage:
				resultMsg += "триггер на каждое сообщение,\n"
			case entities.ScheduleType_EveryHour:
				resultMsg += "каждый час,\n"
			case entities.ScheduleType_EveryDay:
				resultMsg += "каждый день в 12:00,\n"
			}
			var tg string
			if t.Tags != nil {
				tg = *t.Tags
			} else {
				tg = "*нет тегов*"
			}
			resultMsg += "Теги: " + tg + "\n"
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      userId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        resultMsg,
			ReplyMarkup: kb,
		})
	}
}

func nameCheck(name string) bool {
	return tgNameCheck.MatchString(name)
}

func CallBackHandlerReturn(stateManger state_manager.IStateManager) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.NewStartupInlineKeyboard(),
		}

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      userId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Привет! Выбери одно из доступных действий:",
			ReplyMarkup: kb,
		})

		stateManger.ClearState(userId)
	}
}
