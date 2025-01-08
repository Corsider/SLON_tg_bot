package redis

import (
	"SLON_tg_bot/src/domains/entities"
	"SLON_tg_bot/src/state_manager"
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
)

type stateManager struct {
	redisClient *redis.Client
	ctx         context.Context
}

func NewStateManager(redisAddr, redisPass string, redisDb int) state_manager.IStateManager {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: redisPass,
		DB:       redisDb,
	})
	return &stateManager{
		redisClient: client,
		ctx:         context.Background(),
	}
}

func (sm *stateManager) SetState(userID int64, state entities.StateType) {
	key := fmt.Sprintf("user:%d:state", userID)
	sm.redisClient.Set(sm.ctx, key, int(state), 0)
}

func (sm *stateManager) GetState(userID int64) (entities.StateType, bool) {
	key := fmt.Sprintf("user:%d:state", userID)
	result, err := sm.redisClient.Get(sm.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return 0, false
	} else if err != nil {
		panic(err)
	}

	stateInt, err := strconv.Atoi(result)
	if err != nil {
		panic(err)
	}

	return entities.StateType(stateInt), true
}

func (sm *stateManager) ClearState(userID int64) {
	key := fmt.Sprintf("user:%d:state", userID)
	sm.redisClient.Del(sm.ctx, key)
}

func (sm *stateManager) SetUser(userID int64, target string) {
	key := fmt.Sprintf("user:%d:target", userID)
	sm.redisClient.Set(sm.ctx, key, target, 0)
}

func (sm *stateManager) GetUser(userID int64) (string, bool) {
	key := fmt.Sprintf("user:%d:target", userID)
	result, err := sm.redisClient.Get(sm.ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return "", false
	} else if err != nil {
		panic(err)
	}

	return result, true
}
