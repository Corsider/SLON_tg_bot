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
			{Text: "Случайный триггер на сообщение", CallbackData: "sched0"},
		},
		{
			{Text: "Триггер каждые 3 часа", CallbackData: "sched1"},
		},
		{
			{Text: "Триггер на случайного юзера", CallbackData: "sched2"},
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
