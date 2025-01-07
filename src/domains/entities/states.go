package entities

type StateType int

const (
	StateType_Default                    StateType = 0
	StateType_WaitingForTargetName       StateType = 1
	StateType_WaitingForTargetNameToEdit StateType = 2
)
