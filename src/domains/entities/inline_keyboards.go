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
		{
			{Text: "Назад в меню", CallbackData: "return"},
		},
	}
}

func SchedulesInlineKeyboard() [][]models.InlineKeyboardButton {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Триггер на каждое сообщение юзера", CallbackData: "sched0"},
		},
		{
			{Text: "Триггер каждый час", CallbackData: "sched1"},
		},
		{
			{Text: "Триггер каждый день в 12:00", CallbackData: "sched2"},
		},
		{
			{Text: "Назад в меню", CallbackData: "return"},
		},
	}
}

func ReturnInlineKeyboard() [][]models.InlineKeyboardButton {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Отмена", CallbackData: "return"},
		},
	}
}
