package app

import (
	"SLON_tg_bot/src/domains/repositories"
	"SLON_tg_bot/src/state_manager"
	"SLON_tg_bot/src/state_manager/in_memory"
)

type Resources struct {
	StateManager state_manager.IStateManager
	Repository   repositories.IRepository
}

func NewResources() *Resources {
	r := &Resources{}
	r.StateManager = in_memory.NewStateManager()
	// TODO ADD REPO
	r.Repository = repositories.NewEmptyRepository()
	return r
}
