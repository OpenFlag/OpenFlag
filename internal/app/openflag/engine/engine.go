package engine

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/patrickmn/go-cache"
	"github.com/robfig/cron"
	"github.com/sirupsen/logrus"
)

const (
	cacheKey = "flags"
)

type (
	// Evaluation represents one of an evaluation for a flag.
	Evaluation struct {
		Flag    string        `json:"flag"`
		Variant model.Variant `json:"variant"`
	}

	// Result represents of all evaluations for an entity.
	Result struct {
		Entity      model.Entity `json:"entity"`
		Evaluations []Evaluation `json:"evaluations"`
		Timestamp   time.Time    `json:"timestamp"`
	}

	flagSegment struct {
		variant    model.Variant
		constraint constraint.Constraint
	}

	flagItem struct {
		segments []flagSegment
	}
)

// Engine represents an engine interface for the evaluation of an entity.
type Engine interface {
	Evaluate(flags []string, entity model.Entity) (*Result, error)
}

// EvaluationEngine represents an engine for evaluation of an entity.
type EvaluationEngine struct {
	Logger   Logger
	FlagRepo model.FlagRepo
	cache    *cache.Cache
}

// New creates a new evaluation engine.
func New(logger Logger, flagRepo model.FlagRepo) *EvaluationEngine {
	return &EvaluationEngine{
		Logger:   logger,
		FlagRepo: flagRepo,
		cache:    cache.New(cache.NoExpiration, cache.NoExpiration),
	}
}

// Fetch fetches all flags from the database and prepares new period evaluations.
func (e *EvaluationEngine) Fetch() error {
	dbFlags, err := e.FlagRepo.FindAll()
	if err != nil {
		return err
	}

	parse := constraint.Parser{}

	flagMap := map[string]flagItem{}

	for _, dbFlag := range dbFlags {
		flagItem := flagItem{}

		var segments []model.Segment

		if err := json.Unmarshal([]byte(dbFlag.Segments), &segments); err != nil {
			logrus.Errorf(
				"failed to unmarshal db flag %s with id %d into go struct: %s",
				dbFlag.Flag, dbFlag.ID, err.Error(),
			)

			continue
		}

		for _, segment := range segments {
			pco, err := parse.Parse(segment.Expression, segment.Constraints)
			if err != nil {
				logrus.Errorf(
					"failed to parse flag segment expression for flag %s with id %d: %s",
					dbFlag.Flag, dbFlag.ID, err.Error(),
				)

				continue
			}

			co, err := constraint.New(pco.Name, pco.Parameters)
			if err != nil {
				logrus.Errorf(
					"failed to create flag segment constraint for flag %s with id %d: %s",
					dbFlag.Flag, dbFlag.ID, err.Error(),
				)

				continue
			}

			flagItem.segments = append(flagItem.segments, flagSegment{
				variant:    segment.Variant,
				constraint: co,
			})
		}

		flagMap[dbFlag.Flag] = flagItem
	}

	e.cache.Set(cacheKey, flagMap, cache.NoExpiration)

	return nil
}

// Start starts fetching flags from database in periods using given cron pattern.
func (e *EvaluationEngine) Start(cronPattern string) error {
	c := cron.New()

	err := c.AddFunc(cronPattern, func() {
		if err := e.Fetch(); err != nil {
			logrus.Errorf("failed to update flags in period: %s", err.Error())
		}
	})
	if err != nil {
		return err
	}

	c.Start()

	return nil
}

// Evaluate evaluates the given entity.
func (e *EvaluationEngine) Evaluate(flags []string, entity model.Entity) (*Result, error) {
	value, ok := e.cache.Get(cacheKey)
	if !ok {
		return nil, errors.New("failed to load flags from cache")
	}

	flagMap := value.(map[string]flagItem)

	result := Result{
		Entity:      entity,
		Evaluations: []Evaluation{},
		Timestamp:   time.Now(),
	}

	// Use all flags if we don't receive any entered flags.
	if len(flags) == 0 {
		flags = []string{}

		for k := range flagMap {
			flags = append(flags, k)
		}
	}

	for _, flag := range flags {
		f, ok := flagMap[flag]
		if !ok {
			logrus.Warnf("failed to find flag %s in our flags for evaluation", flag)

			continue
		}

		for _, segment := range f.segments {
			if segment.constraint.Evaluate(entity) {
				result.Evaluations = append(result.Evaluations, Evaluation{
					Flag:    flag,
					Variant: segment.variant,
				})

				break
			}
		}
	}

	e.Logger.Log(result)

	return &result, nil
}
