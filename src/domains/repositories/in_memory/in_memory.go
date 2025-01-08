package in_memory

import (
	"SLON_tg_bot/src/domains/entities"
	"SLON_tg_bot/src/domains/repositories"
	"errors"
	"slices"
)

type InMemoryStorage struct {
	users []*entities.TargetUser
}

func NewInMemoryStorage() repositories.IRepository {
	return &InMemoryStorage{
		users: make([]*entities.TargetUser, 0, 10),
	}
}

func (i *InMemoryStorage) AddUser(user *entities.TargetUser) error {
	i.users = append(i.users, user)
	return nil
}

func (i *InMemoryStorage) GetUsersByCreator(creator int64) ([]*entities.TargetUser, error) {
	res := []*entities.TargetUser{}
	for _, u := range i.users {
		if u.CreatorID == creator {
			res = append(res, u)
		}
	}
	return res, nil
}

func (i *InMemoryStorage) GetSingleByCreatorAndTarget(creator int64, target string) (*entities.TargetUser, error) {
	for _, u := range i.users {
		if u.CreatorID == creator && u.Target == target {
			return u, nil
		}
	}
	return nil, errors.New("not found")
}

func (i *InMemoryStorage) RemoveUser(creator int64, target string) error {
	i.users = slices.DeleteFunc(i.users, func(user *entities.TargetUser) bool {
		return user.CreatorID == creator && user.Target == target
	})
	return nil
}

func (i *InMemoryStorage) UpdateUserTags(creator int64, target string, tags string) error {
	u, err := i.GetSingleByCreatorAndTarget(creator, target)
	if err != nil {
		return err
	}
	u.Tags = &tags
	return nil
}

func (i *InMemoryStorage) UpdateUserSched(creator int64, target string, schedType entities.ScheduleType) error {
	u, err := i.GetSingleByCreatorAndTarget(creator, target)
	if err != nil {
		return err
	}
	u.Schedule = schedType
	return nil
}
