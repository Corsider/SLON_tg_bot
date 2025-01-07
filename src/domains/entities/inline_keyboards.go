package entities

import "github.com/go-telegram/bot/models"

func NewStartupInlineKeyboard() [][]models.InlineKeyboardButton {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Добавить юзера", CallbackData: "addUser"},
			{Text: "Редакт. юзера", CallbackData: "editUser"},
		},
		{
			{Text: "Список юзеров", CallbackData: "myUsers"},
		},
	}
}

func EditInlineKeyboard() [][]models.InlineKeyboardButton {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Изм. расписание", CallbackData: "editSched"},
		},
		{
			{Text: "Изм. теги", CallbackData: "editTags"},
		},
		{
			{Text: "Удалить", CallbackData: "editDelete"},
		},
	}
}
