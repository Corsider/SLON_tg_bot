package app

import (
	"SLON_tg_bot/src/domains/repositories"
	"SLON_tg_bot/src/domains/repositories/postgres"
	"SLON_tg_bot/src/state_manager"
	"SLON_tg_bot/src/state_manager/redis"
)

type Resources struct {
	StateManager state_manager.IStateManager
	Repository   repositories.IRepository
}

func NewResources(psqlConnStr, redisConnStr string) (*Resources, error) {
	r := &Resources{}
	r.StateManager = redis.NewStateManager(redisConnStr, "", 0)
	//r.StateManager = in_memory.NewStateManager()
	repo, err := postgres.NewPostgresRepository(psqlConnStr)
	if err != nil {
		return nil, err
	}
	r.Repository = repo
	return r, nil
}
