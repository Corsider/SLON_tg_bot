package repositories

import (
	"SLON_tg_bot/src/domains/entities"
)

type EmptyRepository struct{}

func NewEmptyRepository() IRepository {
	return &EmptyRepository{}
}

func (e *EmptyRepository) AddUser(user *entities.TargetUser) error {
	return nil
}

func (e *EmptyRepository) GetUsersByCreator(creator int64) ([]*entities.TargetUser, error) {
	return nil, nil
}

func (e *EmptyRepository) RemoveUser(creator int64, target string) error {
	return nil
}

func (e *EmptyRepository) GetSingleByCreatorAndTarget(creator int64, target string) (*entities.TargetUser, error) {
	return nil, nil
}

func (e *EmptyRepository) UpdateUserTags(creator int64, target string, tags string) error {
	return nil
}

func (e *EmptyRepository) UpdateUserSched(creator int64, target string, schedType entities.ScheduleType) error {
	return nil
}
