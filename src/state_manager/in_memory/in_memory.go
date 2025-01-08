package in_memory

import (
	"SLON_tg_bot/src/domains/entities"
	"sync"
)

type StateManager struct {
	mu           sync.Mutex
	states       map[int64]entities.StateType
	selectedUser map[int64]string
}

func NewStateManager() *StateManager {
	return &StateManager{
		states:       make(map[int64]entities.StateType),
		selectedUser: make(map[int64]string),
	}
}

func (s *StateManager) SetState(userID int64, state entities.StateType) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.states[userID] = state
}

func (s *StateManager) GetState(userID int64) (entities.StateType, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	state, exists := s.states[userID]
	return state, exists
}

func (s *StateManager) ClearState(userID int64) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.states, userID)
}

func (s *StateManager) SetUser(userID int64, target string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.selectedUser[userID] = target
}

func (s *StateManager) GetUser(userID int64) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	usr, exists := s.selectedUser[userID]
	return usr, exists
}
