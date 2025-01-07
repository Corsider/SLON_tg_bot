package state_manager

import "SLON_tg_bot/src/domains/entities"

type IStateManager interface {
	SetState(userID int64, state entities.StateType)
	GetState(userID int64) (entities.StateType, bool)
	ClearState(userID int64)
}
