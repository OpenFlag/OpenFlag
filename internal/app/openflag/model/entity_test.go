package model_test

import (
	"testing"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/pkg/redis"

	"github.com/stretchr/testify/suite"
)

type EntityRepoSuite struct {
	suite.Suite
	entityRepo model.RedisEntityRepo
}

func (suite *EntityRepoSuite) SetupSuite() {
	cfg := config.Init()

	redisCfg := cfg.Redis

	rMaster, _ := redis.Create(redisCfg.MasterAddress, redisCfg.Options, true)
	rSlave, _ := redis.Create(redisCfg.SlaveAddress, redisCfg.Options, false)

	suite.entityRepo = model.NewRedisEntityRepo(rMaster, rSlave, 1*time.Hour)
}

func (suite *EntityRepoSuite) SetupTest() {
	suite.entityRepo.RedisMaster.FlushDB()
}

func (suite *EntityRepoSuite) TearDownTest() {
	suite.entityRepo.RedisMaster.FlushDB()
}

func (suite *EntityRepoSuite) TestEntityRepo() {
	cases := []struct {
		name             string
		EntitiesForSave  []model.Entity
		EntitiesForQuery []model.Entity
		Result           []model.Entity
	}{
		{
			name: "test case 1",
			EntitiesForSave: []model.Entity{
				{
					EntityID:   1,
					EntityType: "t1",
					EntityContext: map[string]string{
						"k1": "v1",
					},
				},
				{
					EntityID:   2,
					EntityType: "t2",
					EntityContext: map[string]string{
						"k1": "v1",
						"k2": "v2",
					},
				},
			},
			EntitiesForQuery: []model.Entity{
				{
					EntityID:   1,
					EntityType: "t1",
				},
				{
					EntityID:   2,
					EntityType: "t2",
					EntityContext: map[string]string{
						"k2": "v2o",
					},
				},
			},
			Result: []model.Entity{
				{
					EntityID:   1,
					EntityType: "t1",
					EntityContext: map[string]string{
						"k1": "v1",
					},
				},
				{
					EntityID:   2,
					EntityType: "t2",
					EntityContext: map[string]string{
						"k1": "v1",
						"k2": "v2o",
					},
				},
			},
		},
		{
			name: "test case 2",
			EntitiesForSave: []model.Entity{
				{
					EntityID:   1,
					EntityType: "t1",
					EntityContext: map[string]string{
						"k1": "v1",
					},
				},
			},
			EntitiesForQuery: []model.Entity{
				{
					EntityID:   1,
					EntityType: "t1",
				},
				{
					EntityID:   2,
					EntityType: "t2",
				},
			},
			Result: []model.Entity{
				{
					EntityID:   1,
					EntityType: "t1",
					EntityContext: map[string]string{
						"k1": "v1",
					},
				},
				{
					EntityID:   2,
					EntityType: "t2",
				},
			},
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.entityRepo.RedisMaster.FlushDB()

			err := suite.entityRepo.Save(tc.EntitiesForSave)
			suite.NoError(err)

			time.Sleep(2 * time.Second)

			entities, err := suite.entityRepo.Find(tc.EntitiesForQuery)
			suite.NoError(err)

			for i, e := range tc.Result {
				suite.Equal(e, entities[i])
			}
		})
	}
}

func TestEntityRepoSuite(t *testing.T) {
	suite.Run(t, new(EntityRepoSuite))
}
