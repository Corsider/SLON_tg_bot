package entities

type StateType int

const (
	StateType_WaitingForTargetName       StateType = 1
	StateType_WaitingForTargetNameToEdit StateType = 2
	StateType_WaitingForTags             StateType = 3
)
