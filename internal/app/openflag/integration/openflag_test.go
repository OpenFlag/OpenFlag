package integration_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/response"

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
	contextSavingInterval       = 2 * time.Second
)

type OpenFlagSuite struct {
	suite.Suite
	dbClient    *gorm.DB
	redisClient goredis.Cmdable
}

func (suite *OpenFlagSuite) SetupSuite() {
	suite.NoError(os.Setenv("OPENFLAG_SERVER_ADDRESS", serverHTTPAddress))
	suite.NoError(os.Setenv("OPENFLAG_SERVER_GRPC_ADDRESS", serverGRPCAddress))
	suite.NoError(os.Setenv("OPENFLAG_EVALUATION_UPDATE_FLAGS_CRON_PATTERN", "0/5 * * * * *"))

	cfg := config.Init()

	dbClient, err := database.Create(cfg.Database.Driver, cfg.Database.MasterConnStr, cfg.Database.Options)
	suite.NoError(err)

	redisClient, _ := redis.Create(cfg.Redis.MasterAddress, cfg.Redis.Options, true)

	suite.dbClient = dbClient
	suite.redisClient = redisClient

	suite.startServer()
}

func (suite *OpenFlagSuite) TearDownSuite() {
	suite.NoError(os.Unsetenv("OPENFLAG_SERVER_ADDRESS"))
	suite.NoError(os.Unsetenv("OPENFLAG_SERVER_GRPC_ADDRESS"))
	suite.NoError(os.Unsetenv("OPENFLAG_EVALUATION_UPDATE_FLAGS_CRON_PATTERN"))
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
	flagRequest1 := request.CreateFlagRequest{
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

	flagRequest2 := request.CreateFlagRequest{
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
		flagRequestData1, err := json.Marshal(&flagRequest1)
		suite.NoError(err)

		flagRequestData2, err := json.Marshal(&flagRequest2)
		suite.NoError(err)

		req1, err := http.NewRequest("POST", makeURL("/api/v1/flag"), bytes.NewBuffer(flagRequestData1))
		suite.NoError(err)
		req1.Header.Set("Content-Type", "application/json")

		req2, err := http.NewRequest("POST", makeURL("/api/v1/flag"), bytes.NewBuffer(flagRequestData2))
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

	// Wait for syncing flags with database
	time.Sleep(waitingTimeForDoingSomeJobs)

	suite.Run("evaluating entities phase 1", func() {
		evaluationRequest := request.EvaluationRequest{
			Entities: []request.Entity{
				{
					EntityID:   8,
					EntityType: "type1",
				},
				{
					EntityID:   15,
					EntityType: "type2",
				},
				{
					EntityID:   110,
					EntityType: "type2",
				},
			},
			Flags: []string{"flag1", "flag2"},
		}

		evaluationRequestData, err := json.Marshal(&evaluationRequest)
		suite.NoError(err)

		req, err := http.NewRequest("POST", makeURL("/api/v1/evaluation"), bytes.NewBuffer(evaluationRequestData))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}

		resp, err := client.Do(req)
		suite.NoError(err)
		suite.Equal(http.StatusOK, resp.StatusCode)

		respBody, err := ioutil.ReadAll(resp.Body)
		suite.NoError(err)

		suite.NoError(resp.Body.Close())

		var respParsed []response.EvaluationResponse

		suite.NoError(json.Unmarshal(respBody, &respParsed))

		expectedResponse := []response.EvaluationResponse{
			{
				Entity: response.Entity{
					EntityID:   8,
					EntityType: "type1",
				},
				Evaluations: []response.Evaluation{
					{
						Flag: "flag1",
						Variant: response.Variant{
							VariantKey:        "on1",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
					{
						Flag: "flag2",
						Variant: response.Variant{
							VariantKey:        "on2",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
				},
			},
			{
				Entity: response.Entity{
					EntityID:   15,
					EntityType: "type2",
				},
				Evaluations: []response.Evaluation{
					{
						Flag: "flag1",
						Variant: response.Variant{
							VariantKey:        "off1",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
					{
						Flag: "flag2",
						Variant: response.Variant{
							VariantKey:        "on2",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
				},
			},
			{
				Entity: response.Entity{
					EntityID:   110,
					EntityType: "type2",
				},
				Evaluations: []response.Evaluation{
					{
						Flag: "flag1",
						Variant: response.Variant{
							VariantKey:        "off1",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
					{
						Flag: "flag2",
						Variant: response.Variant{
							VariantKey:        "off2",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
				},
			},
		}

		suite.Equal(expectedResponse, respParsed)
	})

	suite.Run("evaluating entities phase 2", func() {
		evaluationRequest := request.EvaluationRequest{
			Entities: []request.Entity{
				{
					EntityID:   8,
					EntityType: "type1",
					EntityContext: map[string]string{
						"c1": "v1",
					},
				},
			},
			Flags:        []string{"flag1", "flag2"},
			SaveContexts: true,
		}

		evaluationRequestData, err := json.Marshal(&evaluationRequest)
		suite.NoError(err)

		req, err := http.NewRequest("POST", makeURL("/api/v1/evaluation"), bytes.NewBuffer(evaluationRequestData))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}

		resp, err := client.Do(req)
		suite.NoError(err)
		suite.Equal(http.StatusOK, resp.StatusCode)

		respBody, err := ioutil.ReadAll(resp.Body)
		suite.NoError(err)

		suite.NoError(resp.Body.Close())

		var respParsed []response.EvaluationResponse

		suite.NoError(json.Unmarshal(respBody, &respParsed))

		expectedResponse := []response.EvaluationResponse{
			{
				Entity: response.Entity{
					EntityID:   8,
					EntityType: "type1",
					EntityContext: map[string]string{
						"c1": "v1",
					},
				},
				Evaluations: []response.Evaluation{
					{
						Flag: "flag1",
						Variant: response.Variant{
							VariantKey:        "on1",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
					{
						Flag: "flag2",
						Variant: response.Variant{
							VariantKey:        "on2",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
				},
			},
		}

		suite.Equal(expectedResponse, respParsed)
	})

	// Wait for saving contexts to Redis
	time.Sleep(contextSavingInterval)

	suite.Run("evaluating entities phase 3", func() {
		evaluationRequest := request.EvaluationRequest{
			Entities: []request.Entity{
				{
					EntityID:   8,
					EntityType: "type1",
					EntityContext: map[string]string{
						"c2": "v2",
					},
				},
			},
			Flags:             []string{"flag1", "flag2"},
			SaveContexts:      true,
			UseStoredContexts: true,
		}

		evaluationRequestData, err := json.Marshal(&evaluationRequest)
		suite.NoError(err)

		req, err := http.NewRequest("POST", makeURL("/api/v1/evaluation"), bytes.NewBuffer(evaluationRequestData))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}

		resp, err := client.Do(req)
		suite.NoError(err)
		suite.Equal(http.StatusOK, resp.StatusCode)

		respBody, err := ioutil.ReadAll(resp.Body)
		suite.NoError(err)

		suite.NoError(resp.Body.Close())

		var respParsed []response.EvaluationResponse

		suite.NoError(json.Unmarshal(respBody, &respParsed))

		expectedResponse := []response.EvaluationResponse{
			{
				Entity: response.Entity{
					EntityID:   8,
					EntityType: "type1",
					EntityContext: map[string]string{
						"c1": "v1",
						"c2": "v2",
					},
				},
				Evaluations: []response.Evaluation{
					{
						Flag: "flag1",
						Variant: response.Variant{
							VariantKey:        "on1",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
					{
						Flag: "flag2",
						Variant: response.Variant{
							VariantKey:        "on2",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
				},
			},
		}

		suite.Equal(expectedResponse, respParsed)
	})

	// Wait for saving contexts to Redis
	time.Sleep(contextSavingInterval)

	suite.Run("evaluating entities phase 4", func() {
		evaluationRequest := request.EvaluationRequest{
			Entities: []request.Entity{
				{
					EntityID:   8,
					EntityType: "type1",
				},
			},
			Flags:             []string{"flag1", "flag2"},
			SaveContexts:      false,
			UseStoredContexts: true,
		}

		evaluationRequestData, err := json.Marshal(&evaluationRequest)
		suite.NoError(err)

		req, err := http.NewRequest("POST", makeURL("/api/v1/evaluation"), bytes.NewBuffer(evaluationRequestData))
		suite.NoError(err)
		req.Header.Set("Content-Type", "application/json")

		client := http.Client{}

		resp, err := client.Do(req)
		suite.NoError(err)
		suite.Equal(http.StatusOK, resp.StatusCode)

		respBody, err := ioutil.ReadAll(resp.Body)
		suite.NoError(err)

		suite.NoError(resp.Body.Close())

		var respParsed []response.EvaluationResponse

		suite.NoError(json.Unmarshal(respBody, &respParsed))

		expectedResponse := []response.EvaluationResponse{
			{
				Entity: response.Entity{
					EntityID:   8,
					EntityType: "type1",
					EntityContext: map[string]string{
						"c1": "v1",
						"c2": "v2",
					},
				},
				Evaluations: []response.Evaluation{
					{
						Flag: "flag1",
						Variant: response.Variant{
							VariantKey:        "on1",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
					{
						Flag: "flag2",
						Variant: response.Variant{
							VariantKey:        "on2",
							VariantAttachment: json.RawMessage(`{}`),
						},
					},
				},
			},
		}

		suite.Equal(expectedResponse, respParsed)
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
