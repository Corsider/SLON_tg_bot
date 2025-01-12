package handlers

import (
	"SLON_tg_bot/src/domains/entities"
	"SLON_tg_bot/src/domains/repositories"
	"SLON_tg_bot/src/state_manager"
	"context"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"log"
	"slices"
)

func CallBackHandlerEditTags(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.AddTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.ToFlatTags()

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
				"\n\nСейчас юзеру назначены теги:\n" + currentTags,
			ReplyMarkup: kb,
		})

		stateManager.ClearState(userId)
	}
}

func CallBackHandlerDelTags(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.DelTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.ToFlatTags()

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Выбери теги для удаления. По крайней мере 1 тег должен остаться назначенным." +
				"\n\nСейчас юзеру назначены теги:\n" + currentTags,
			ReplyMarkup: kb,
		})

		stateManager.ClearState(userId)
	}
}

func CallBackHandlerTagReady(stateManager state_manager.IStateManager) func(ctx context.Context, b *bot.Bot, update *models.Update) {
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
			log.Printf("[APP] [ERR] Error occured, user does not exist.")
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:      userId,
			MessageID:   update.CallbackQuery.Message.Message.ID,
			Text:        "Теги пользователя " + target + " обновлены.\n" + spravka,
			ReplyMarkup: kb,
		})

		stateManager.ClearState(userId)
	}
}

func CallBackHandlerDelTag0(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.DelTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.GetTags()
		if len(currentTags) == 1 {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "Нельзя удалить тег, т.к. у юзера должен быть хотя бы 1 тег." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})
			stateManager.ClearState(userId)
			return
		}
		if !slices.Contains(currentTags, entities.TagType_INSULT) {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "У юзера нет этого тега." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})
			stateManager.ClearState(userId)
			return
		}

		currentTags = slices.DeleteFunc(currentTags, func(tagType entities.TagType) bool {
			return tagType == entities.TagType_INSULT
		})

		err = repo.UpdateUserTags(userId, target, currentTags)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		u, err = repo.GetSingleByCreatorAndTarget(userId, target)

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Готово." +
				"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
			ReplyMarkup: kb,
		})
		log.Printf("[APP] [INFO] User %d has updated target's (%s) tags to %s", userId, target, u.ToFlatTags())
		stateManager.ClearState(userId)
	}
}

func CallBackHandlerDelTag1(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.DelTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.GetTags()
		if len(currentTags) == 1 {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "Нельзя удалить тег, т.к. у юзера должен быть хотя бы 1 тег." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})
			stateManager.ClearState(userId)
			return
		}
		if !slices.Contains(currentTags, entities.TagType_OBSCENITY) {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "У юзера нет этого тега." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})
			stateManager.ClearState(userId)
			return
		}

		currentTags = slices.DeleteFunc(currentTags, func(tagType entities.TagType) bool {
			return tagType == entities.TagType_OBSCENITY
		})

		err = repo.UpdateUserTags(userId, target, currentTags)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		u, err = repo.GetSingleByCreatorAndTarget(userId, target)

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Готово." +
				"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
			ReplyMarkup: kb,
		})
		log.Printf("[APP] [INFO] User %d has updated target's (%s) tags to %s", userId, target, u.ToFlatTags())
		stateManager.ClearState(userId)
	}
}

func CallBackHandlerDelTag2(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.DelTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.GetTags()
		if len(currentTags) == 1 {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "Нельзя удалить тег, т.к. у юзера должен быть хотя бы 1 тег." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})
			stateManager.ClearState(userId)
			return
		}
		if !slices.Contains(currentTags, entities.TagType_THREAT) {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "У юзера нет этого тега." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})
			stateManager.ClearState(userId)
			return
		}

		currentTags = slices.DeleteFunc(currentTags, func(tagType entities.TagType) bool {
			return tagType == entities.TagType_THREAT
		})

		err = repo.UpdateUserTags(userId, target, currentTags)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		u, err = repo.GetSingleByCreatorAndTarget(userId, target)

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Готово." +
				"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
			ReplyMarkup: kb,
		})
		log.Printf("[APP] [INFO] User %d has updated target's (%s) tags to %s", userId, target, u.ToFlatTags())
		stateManager.ClearState(userId)
	}
}

func CallBackHandlerAssignTag0(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.AddTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.GetTags()
		if slices.Contains(currentTags, entities.TagType_INSULT) {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})

			stateManager.ClearState(userId)
			return
		}

		err = repo.UpdateUserTags(userId, target, append(u.GetTags(), entities.TagType_INSULT))
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		u, err = repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
				"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
			ReplyMarkup: kb,
		})
		log.Printf("[APP] [INFO] User %d has updated target's (%s) tags to %s", userId, target, u.ToFlatTags())
		stateManager.ClearState(userId)
	}
}

func CallBackHandlerAssignTag1(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.AddTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.GetTags()
		if slices.Contains(currentTags, entities.TagType_OBSCENITY) {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})

			stateManager.ClearState(userId)
			return
		}

		err = repo.UpdateUserTags(userId, target, append(u.GetTags(), entities.TagType_OBSCENITY))
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		u, err = repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
				"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
			ReplyMarkup: kb,
		})
		log.Printf("[APP] [INFO] User %d has updated target's (%s) tags to %s", userId, target, u.ToFlatTags())
		stateManager.ClearState(userId)
	}
}

func CallBackHandlerAssignTag2(stateManager state_manager.IStateManager, repo repositories.IRepository) func(ctx context.Context, b *bot.Bot, update *models.Update) {
	return func(ctx context.Context, b *bot.Bot, update *models.Update) {
		userId := update.CallbackQuery.Message.Message.Chat.ID

		b.AnswerCallbackQuery(ctx, &bot.AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
			ShowAlert:       false,
		})

		kb := &models.InlineKeyboardMarkup{
			InlineKeyboard: entities.AddTagsInlineKeyboard(),
		}

		target, exists := stateManager.GetUser(userId)
		if !exists {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Target not found for user %d.", userId)
			return
		}

		u, err := repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		currentTags := u.GetTags()
		if slices.Contains(currentTags, entities.TagType_THREAT) {
			b.EditMessageText(ctx, &bot.EditMessageTextParams{
				ChatID:    userId,
				MessageID: update.CallbackQuery.Message.Message.ID,
				Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
					"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
				ReplyMarkup: kb,
			})

			stateManager.ClearState(userId)
			return
		}

		err = repo.UpdateUserTags(userId, target, append(u.GetTags(), entities.TagType_THREAT))
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		u, err = repo.GetSingleByCreatorAndTarget(userId, target)
		if err != nil {
			errorMessage(ctx, b, userId)
			log.Printf("[APP] [ERR] Error occured: %v", err)
			return
		}

		b.EditMessageText(ctx, &bot.EditMessageTextParams{
			ChatID:    userId,
			MessageID: update.CallbackQuery.Message.Message.ID,
			Text: "Выбери теги для уточнения поведения модели. Можно выбрать от 1 до 3 тегов." +
				"\n\nСейчас юзеру назначены теги:\n" + u.ToFlatTags(),
			ReplyMarkup: kb,
		})
		log.Printf("[APP] [INFO] User %d has updated target's (%s) tags to %s", userId, target, u.ToFlatTags())
		stateManager.ClearState(userId)
	}
}
