package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

const (
	entityName           = "redis_entity"
	pipelineExecDuration = 1 * time.Second
)

type (
	// Entity represents the context of what we are going to assign the variant on.
	// Usually, OpenFlag expects the context coming with the entity,
	// so that one can define constraints based on the context of the entity.
	Entity struct {
		EntityID      int64             `json:"entity_id"`
		EntityType    string            `json:"entity_type"`
		EntityContext map[string]string `json:"entity_context,omitempty"`
	}
)

// EntityRepo represents an interface for working with persist entities.
type EntityRepo interface {
	Save(entities []Entity) error
	Find(entities []Entity) ([]Entity, error)
}

// RedisEntityRepo is an implementation of EntityRepo for Redis.
type RedisEntityRepo struct {
	RedisMaster redis.Cmdable
	RedisSlave  redis.Cmdable
	Pipeliner   redis.Pipeliner
	Expiration  time.Duration
}

// NewRedisEntityRepo created a new Redis entity repo.
func NewRedisEntityRepo(rMaster redis.Cmdable, rSlave redis.Cmdable, expiration time.Duration) RedisEntityRepo {
	repo := RedisEntityRepo{
		RedisMaster: rMaster,
		RedisSlave:  rSlave,
		Pipeliner:   rMaster.Pipeline(),
		Expiration:  expiration,
	}

	repo.runPipeline()

	return repo
}

func (r RedisEntityRepo) runPipeline() {
	go func() {
		timer := time.NewTicker(pipelineExecDuration).C
		for range timer {
			startTime := time.Now()

			_, err := r.Pipeliner.Exec()

			metrics.report(entityName, "run_pipeline", startTime, err)
		}
	}()
}

func (r RedisEntityRepo) entityIDKey(eID int64, eType string) string {
	return fmt.Sprintf("openflag:entity:type:%s:id:%d", eType, eID)
}

// Save saves entity context to Redis.
func (r RedisEntityRepo) Save(entities []Entity) (finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(entityName, "save", startTime, finalErr) }()

	for _, e := range entities {
		rKey := r.entityIDKey(e.EntityID, e.EntityType)

		context, err := json.Marshal(e.EntityContext)
		if err != nil {
			return err
		}

		r.Pipeliner.Set(rKey, context, r.Expiration)
	}

	return nil
}

// Find finds entity context from Redis.
func (r RedisEntityRepo) Find(entities []Entity) (_ []Entity, finalErr error) {
	startTime := time.Now()

	defer func() { metrics.report(entityName, "find", startTime, finalErr) }()

	result, err := r.RedisSlave.Pipelined(func(pipeliner redis.Pipeliner) error {
		for _, entity := range entities {
			pipeliner.Get(r.entityIDKey(entity.EntityID, entity.EntityType))
		}

		return nil
	})
	if err != nil && err != redis.Nil {
		return nil, err
	}

	for i := range result {
		redisContext, ok := result[i].(*redis.StringCmd)
		if !ok {
			return nil, errors.New("cannot read pipeline result")
		}

		jsonContext, err := redisContext.Result()
		if err != nil && err != redis.Nil {
			return nil, err
		} else if err == redis.Nil {
			continue
		}

		context := map[string]string{}

		err = json.Unmarshal([]byte(jsonContext), &context)
		if err != nil {
			return nil, err
		}

		for key, value := range context {
			_, ok := entities[i].EntityContext[key]
			if ok {
				continue
			}

			if entities[i].EntityContext == nil {
				entities[i].EntityContext = map[string]string{}
			}

			entities[i].EntityContext[key] = value
		}
	}

	return entities, nil
}
