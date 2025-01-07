package repositories

import "SLON_tg_bot/src/domains/entities"

type IRepository interface {
	AddUser(user *entities.TargetUser) error
	GetUsersByCreator(creator int64) ([]*entities.TargetUser, error)
	UpdateUser(creator int64, upd *entities.TargetUser) error
}
