package entities

import "github.com/go-telegram/bot/models"

type StartupInlineKeyboard [][]models.InlineKeyboardButton

func NewStartupInlineKeyboard() StartupInlineKeyboard {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Добавить юзера", CallbackData: "addUser"},
			{Text: "Редактировать юзера", CallbackData: "editUser"},
		},
	}
}
