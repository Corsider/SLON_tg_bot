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
			{Text: "Доб. теги", CallbackData: "addTags"},
		},
		{
			{Text: "Удал. теги", CallbackData: "delTags"},
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

func AddTagsInlineKeyboard() [][]models.InlineKeyboardButton {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Оскорбление", CallbackData: "tag0"},
		},
		{
			{Text: "Непристойность", CallbackData: "tag1"},
		},
		{
			{Text: "Угроза", CallbackData: "tag2"},
		},
		{
			{Text: "Готово", CallbackData: "tagReady"},
		},
	}
}

func DelTagsInlineKeyboard() [][]models.InlineKeyboardButton {
	return [][]models.InlineKeyboardButton{
		{
			{Text: "Оскорбление", CallbackData: "dtag0"},
		},
		{
			{Text: "Непристойность", CallbackData: "dtag1"},
		},
		{
			{Text: "Угроза", CallbackData: "dtag2"},
		},
		{
			{Text: "Готово", CallbackData: "tagReady"},
		},
	}
}
