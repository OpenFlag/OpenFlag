package openflag_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/request"
	"github.com/OpenFlag/OpenFlag/pkg/database"
	"github.com/OpenFlag/OpenFlag/pkg/redis"
	goredis "github.com/go-redis/redis"
	"github.com/jinzhu/gorm"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/stretchr/testify/suite"
)

const (
	serverHTTPAddress           = "127.0.0.1:7677"
	serverGRPCAddress           = "127.0.0.1:7678"
	waitingTimeForDoingSomeJobs = 10 * time.Second
)

type OpenFlagSuite struct {
	suite.Suite
	dbClient    *gorm.DB
	redisClient goredis.Cmdable
}

func (suite *OpenFlagSuite) SetupSuite() {
	cfg := config.Init()

	dbClient, err := database.Create(cfg.Database.Driver, cfg.Database.MasterConnStr, cfg.Database.Options)
	suite.NoError(err)

	redisClient, _ := redis.Create(cfg.Redis.MasterAddress, cfg.Redis.Options, true)

	suite.dbClient = dbClient
	suite.redisClient = redisClient

	suite.NoError(os.Setenv("OPENFLAG_SERVER_ADDRESS", serverHTTPAddress))
	suite.NoError(os.Setenv("OPENFLAG_GRPC_SERVER_ADDRESS", serverGRPCAddress))

	suite.startServer()
}

func (suite *OpenFlagSuite) TearDownSuite() {
	suite.NoError(os.Unsetenv("OPENFLAG_SERVER_ADDRESS"))
	suite.NoError(os.Unsetenv("OPENFLAG_GRPC_SERVER_ADDRESS"))
}

func (suite *OpenFlagSuite) SetupTest() {
	suite.NoError(suite.redisClient.FlushDB().Err())
	suite.NoError(suite.dbClient.Exec(`truncate table flags`).Error)
}

func (suite *OpenFlagSuite) TearDownTest() {
	suite.NoError(suite.redisClient.FlushDB().Err())
	suite.NoError(suite.dbClient.Exec(`truncate table flags`).Error)
}

func (suite *OpenFlagSuite) TestScenario() {
	flag1Request := request.CreateFlagRequest{
		Flag: request.Flag{
			Tags:        []string{"tag1", "tag2"},
			Description: "flag1 description",
			Flag:        "flag1",
			Segments: []request.Segment{
				{
					Description: "segment 1",
					Constraints: map[string]request.Constraint{
						"A": {
							Name:       constraint.LessThanConstraintName,
							Parameters: json.RawMessage(`{"value": 10}`),
						},
						"B": {
							Name:       constraint.BiggerThanConstraintName,
							Parameters: json.RawMessage(`{"value": 5}`),
						},
					},
					Expression: fmt.Sprintf("A %s B", constraint.IntersectionConstraintName),
					Variant: request.Variant{
						VariantKey:        "on1",
						VariantAttachment: json.RawMessage(`{}`),
					},
				},
				{
					Description: "segment 2",
					Constraints: map[string]request.Constraint{
						"A": {
							Name:       constraint.AlwaysConstraintName,
							Parameters: json.RawMessage(`{}`),
						},
					},
					Expression: "A",
					Variant: request.Variant{
						VariantKey:        "off1",
						VariantAttachment: json.RawMessage(`{}`),
					},
				},
			},
		},
	}

	flag2Request := request.CreateFlagRequest{
		Flag: request.Flag{
			Tags:        []string{"tag2", "tag3"},
			Description: "flag2 description",
			Flag:        "flag2",
			Segments: []request.Segment{
				{
					Description: "segment 1",
					Constraints: map[string]request.Constraint{
						"A": {
							Name:       constraint.LessThanConstraintName,
							Parameters: json.RawMessage(`{"value": 100}`),
						},
						"B": {
							Name:       constraint.BiggerThanConstraintName,
							Parameters: json.RawMessage(`{"value": 5}`),
						},
					},
					Expression: fmt.Sprintf("A %s B", constraint.IntersectionConstraintName),
					Variant: request.Variant{
						VariantKey:        "on2",
						VariantAttachment: json.RawMessage(`{}`),
					},
				},
				{
					Description: "segment 2",
					Constraints: map[string]request.Constraint{
						"A": {
							Name:       constraint.AlwaysConstraintName,
							Parameters: json.RawMessage(`{}`),
						},
					},
					Expression: "A",
					Variant: request.Variant{
						VariantKey:        "off2",
						VariantAttachment: json.RawMessage(`{}`),
					},
				},
			},
		},
	}

	suite.Run("create feature flags", func() {
		flag1RequestData, err := json.Marshal(&flag1Request)
		suite.NoError(err)

		flag2RequestData, err := json.Marshal(&flag2Request)
		suite.NoError(err)

		req1, err := http.NewRequest("POST", makeURL("/api/v1/flag"), bytes.NewBuffer(flag1RequestData))
		suite.NoError(err)
		req1.Header.Set("Content-Type", "application/json")

		req2, err := http.NewRequest("POST", makeURL("/api/v1/flag"), bytes.NewBuffer(flag2RequestData))
		suite.NoError(err)
		req2.Header.Set("Content-Type", "application/json")

		client := http.Client{}

		resp1, err := client.Do(req1)
		suite.NoError(err)
		suite.NoError(resp1.Body.Close())
		suite.Equal(http.StatusOK, resp1.StatusCode)

		resp2, err := client.Do(req2)
		suite.NoError(err)
		suite.NoError(resp2.Body.Close())
		suite.Equal(http.StatusOK, resp2.StatusCode)
	})
}

func (suite *OpenFlagSuite) startServer() {
	root := cmd.NewRootCommand()
	if root == nil {
		suite.Fail("root command is nil!")
		return
	}

	root.SetArgs([]string{"server"})

	go func() {
		if err := root.Execute(); err != nil {
			suite.Fail("failed to execute root command: %s", err.Error())
		}
	}()

	// Wait for starting the server
	time.Sleep(waitingTimeForDoingSomeJobs)
}

func makeURL(path string) string {
	return fmt.Sprintf("http://%s%s", serverHTTPAddress, path)
}

func TestOpenFlagSuite(t *testing.T) {
	suite.Run(t, new(OpenFlagSuite))
}
