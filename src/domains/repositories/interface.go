package repositories

import "SLON_tg_bot/src/domains/entities"

type IRepository interface {
	AddUser(user *entities.TargetUser) error
	GetUsersByCreator(creator int64) ([]*entities.TargetUser, error)
	GetSingleByCreatorAndTarget(creator int64, target string) (*entities.TargetUser, error)
	RemoveUser(creator int64, target string) error
	UpdateUserTags(creator int64, target string, tags []entities.TagType) error
	UpdateUserSched(creator int64, target string, schedType entities.ScheduleType) error
}
