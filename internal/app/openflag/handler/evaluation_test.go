package handler_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/engine"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/handler"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/request"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/response"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type (
	fakeEntityRepo struct {
		model.EntityRepo
		saveFunc func(entities []model.Entity) error
		findFunc func(entities []model.Entity) (entities2 []model.Entity, e error)
	}

	fakeEvaluationEngine struct {
		engine.Engine
		evaluateFunc func(flags []string, entity model.Entity) (result *engine.Result, e error)
	}
)

func (f fakeEntityRepo) Save(entities []model.Entity) error {
	return f.saveFunc(entities)
}

func (f fakeEntityRepo) Find(entities []model.Entity) ([]model.Entity, error) {
	return f.findFunc(entities)
}

func (f fakeEvaluationEngine) Evaluate(flags []string, entity model.Entity) (*engine.Result, error) {
	return f.evaluateFunc(flags, entity)
}

type EvaluationHandlerSuite struct {
	suite.Suite
	engine               *echo.Echo
	fakeEntityRepo       *fakeEntityRepo
	fakeEvaluationEngine *fakeEvaluationEngine
}

func (suite *EvaluationHandlerSuite) SetupSuite() {
	suite.engine = echo.New()

	suite.fakeEntityRepo = &fakeEntityRepo{}
	suite.fakeEvaluationEngine = &fakeEvaluationEngine{}

	suite.engine.POST("/v1/evaluation", handler.EvaluationHandler{
		Engine:     suite.fakeEvaluationEngine,
		EntityRepo: suite.fakeEntityRepo,
	}.Evaluate)
}

func (suite *EvaluationHandlerSuite) TestEvaluation() {
	cases := []struct {
		name             string
		req              request.EvaluationRequest
		contextsStore    func(entities []model.Entity) error
		contextsReader   func(entities []model.Entity) ([]model.Entity, error)
		evaluationEngine func(flags []string, entity model.Entity) (*engine.Result, error)
		expectedStatus   int
		expectedResponse []response.EvaluationResponse
	}{
		{
			name: "successfully send evaluation request and get response 1",
			req: request.EvaluationRequest{
				Entities: []request.Entity{
					{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
						},
					},
				},
				Flags:             []string{"flag1"},
				SaveContexts:      false,
				UseStoredContexts: false,
			},
			contextsStore:  func(entities []model.Entity) error { return nil },
			contextsReader: func(entities []model.Entity) (entities2 []model.Entity, e error) { return nil, nil },
			evaluationEngine: func(flags []string, entity model.Entity) (result *engine.Result, e error) {
				return &engine.Result{
					Entity: entity,
					Evaluations: []engine.Evaluation{
						{
							Flag: "flag1",
							Variant: model.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
					Timestamp: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []response.EvaluationResponse{
				{
					Entity: response.Entity{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
						},
					},
					Evaluations: []response.Evaluation{
						{
							Flag: "flag1",
							Variant: response.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
				},
			},
		},
		{
			name: "successfully send evaluation request and get response 2",
			req: request.EvaluationRequest{
				Entities: []request.Entity{
					{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
						},
					},
				},
				Flags:             []string{"flag1"},
				SaveContexts:      true,
				UseStoredContexts: true,
			},
			contextsStore: func(entities []model.Entity) error { return nil },
			contextsReader: func(entities []model.Entity) (entities2 []model.Entity, e error) {
				return []model.Entity{
					{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
					},
				}, nil
			},
			evaluationEngine: func(flags []string, entity model.Entity) (result *engine.Result, e error) {
				return &engine.Result{
					Entity: entity,
					Evaluations: []engine.Evaluation{
						{
							Flag: "flag1",
							Variant: model.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
					Timestamp: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []response.EvaluationResponse{
				{
					Entity: response.Entity{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
					},
					Evaluations: []response.Evaluation{
						{
							Flag: "flag1",
							Variant: response.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
				},
			},
		},
		{
			name: "successfully send evaluation request and get response 3",
			req: request.EvaluationRequest{
				Entities: []request.Entity{
					{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
						},
					},
					{
						EntityID:   2,
						EntityType: "type2",
						EntityContext: map[string]string{
							"k1": "v1",
						},
					},
				},
				Flags:             []string{"flag1"},
				SaveContexts:      true,
				UseStoredContexts: true,
			},
			contextsStore: func(entities []model.Entity) error { return nil },
			contextsReader: func(entities []model.Entity) (entities2 []model.Entity, e error) {
				return []model.Entity{
					{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
					},
					{
						EntityID:   2,
						EntityType: "type2",
						EntityContext: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
					},
				}, nil
			},
			evaluationEngine: func(flags []string, entity model.Entity) (result *engine.Result, e error) {
				return &engine.Result{
					Entity: entity,
					Evaluations: []engine.Evaluation{
						{
							Flag: "flag1",
							Variant: model.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
					Timestamp: time.Now(),
				}, nil
			},
			expectedStatus: http.StatusOK,
			expectedResponse: []response.EvaluationResponse{
				{
					Entity: response.Entity{
						EntityID:   1,
						EntityType: "type1",
						EntityContext: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
					},
					Evaluations: []response.Evaluation{
						{
							Flag: "flag1",
							Variant: response.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
				},
				{
					Entity: response.Entity{
						EntityID:   2,
						EntityType: "type2",
						EntityContext: map[string]string{
							"k1": "v1",
							"k2": "v2",
						},
					},
					Evaluations: []response.Evaluation{
						{
							Flag: "flag1",
							Variant: response.Variant{
								VariantKey:        "on",
								VariantAttachment: json.RawMessage(`{}`),
							},
						},
					},
				},
			},
		},
		{
			name: "failed to send evaluation request and get response 1",
			req: request.EvaluationRequest{
				Entities: []request.Entity{
					{
						EntityID:   1,
						EntityType: "",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "failed to send evaluation request and get response 2",
			req: request.EvaluationRequest{
				Entities: []request.Entity{
					{
						EntityID:   0,
						EntityType: "type1",
					},
				},
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeEntityRepo.saveFunc = tc.contextsStore
			suite.fakeEntityRepo.findFunc = tc.contextsReader
			suite.fakeEvaluationEngine.evaluateFunc = tc.evaluationEngine

			data, err := json.Marshal(tc.req)
			suite.NoError(err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/v1/evaluation", bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.expectedStatus, w.Code, tc.name)

			if tc.expectedStatus == http.StatusOK {
				var resp []response.EvaluationResponse

				suite.NoError(json.Unmarshal(w.Body.Bytes(), &resp))

				suite.Equal(tc.expectedResponse, resp)
			}
		})
	}
}

func TestEvaluationHandlerSuite(t *testing.T) {
	suite.Run(t, new(EvaluationHandlerSuite))
}
