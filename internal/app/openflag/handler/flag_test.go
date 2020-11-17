package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/response"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/constraint"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/handler"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/request"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/suite"
)

type fakeFlagRepo struct {
	model.FlagRepo
	repoError    error
	toBeDeleteID int64
	toBeUpdateID int64
}

func (f *fakeFlagRepo) Create(flag *model.Flag) error {
	return f.repoError
}

func (f *fakeFlagRepo) Delete(id int64) error {
	if f.repoError != nil {
		return f.repoError
	}

	f.toBeDeleteID = id

	return nil
}

func (f *fakeFlagRepo) Update(id int64, flag *model.Flag) error {
	if f.repoError != nil {
		return f.repoError
	}

	f.toBeUpdateID = id

	return nil
}

func (f *fakeFlagRepo) FindByID(id int64) (*model.Flag, error) {
	if f.repoError != nil {
		return nil, f.repoError
	}

	tags := `["tag1"]`

	flag := &model.Flag{
		ID:          10,
		Tags:        &tags,
		Description: "description 1",
		Flag:        "flag1",
		Segments: `
			[
				{
					"description": "segment 1",
					"constraints": {
						"A": {
							"name": "<",
							"parameters": {
								"value": 10
							}
						},
						"B": {
							"name": ">",
							"parameters": {
								"value": 5
							}
						}
					},
					"expression": "A ∩ B",
					"variant": {
						"key": "on"
					}
				}
			]
		`,
	}

	if id == flag.ID {
		return flag, nil
	}

	return nil, model.ErrFlagNotFound
}

type FlagHandlerSuite struct {
	suite.Suite
	engine       *echo.Echo
	fakeFlagRepo *fakeFlagRepo
}

func (suite *FlagHandlerSuite) SetupSuite() {
	suite.fakeFlagRepo = &fakeFlagRepo{}
	suite.engine = echo.New()

	suite.engine.POST("/v1/flag", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.Create)
	suite.engine.DELETE("/v1/flag/:id", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.Delete)
	suite.engine.PUT("/v1/flag/:id", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.Update)
	suite.engine.GET("/v1/flag/:id", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.FindByID)
	suite.engine.POST("/v1/flag/tag", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.FindByTag)
	suite.engine.POST("/v1/flag/history", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.FindByFlag)
	suite.engine.POST("/v1/flags", handler.FlagHandler{FlagRepo: suite.fakeFlagRepo}.FindFlags)
}

func (suite *FlagHandlerSuite) TestCreateFlag() {
	cases := []struct {
		name      string
		req       request.CreateFlagRequest
		status    int
		repoError error
	}{
		{
			name: "successfully create flag 1",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusOK,
		},
		{
			name: "successfully create flag 2",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusOK,
		},
		{
			name: "failed to create flag 1",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: errors.New("fake flag repo error"),
			status:    http.StatusInternalServerError,
		},
		{
			name: "failed to create flag 2",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "failed to create flag 3",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "failed to create flag 4",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "failed to create flag 5",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
							Constraints: map[string]request.Constraint{},
							Expression:  fmt.Sprintf("A %s B", constraint.IntersectionConstraintName),
							Variant: request.Variant{
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "failed to create flag 6",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
							Variant: request.Variant{
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "failed to create flag 7",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
							Variant:    request.Variant{},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "failed to create flag 8",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
				},
			},
			repoError: nil,
			status:    http.StatusBadRequest,
		},
		{
			name: "successfully create flag 9",
			req: request.CreateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: model.ErrDuplicateFlagFound,
			status:    http.StatusConflict,
		},
	}

	for i := range cases {
		tc := cases[i]

		suite.Run(tc.name, func() {
			suite.fakeFlagRepo.repoError = tc.repoError

			data, err := json.Marshal(tc.req)
			suite.NoError(err)

			w := httptest.NewRecorder()
			req := httptest.NewRequest("POST", "/v1/flag", bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code, tc.name)
		})
	}
}

func (suite *FlagHandlerSuite) TestDeleteFlag() {
	cases := []struct {
		name      string
		flagID    string
		status    int
		repoError error
	}{
		{
			name:      "successfully delete flag 1",
			flagID:    "10",
			status:    http.StatusNoContent,
			repoError: nil,
		},
		{
			name:      "failed to delete flag 1",
			flagID:    "10",
			status:    http.StatusInternalServerError,
			repoError: errors.New("fake flag repo error"),
		},
		{
			name:      "failed to delete flag 2",
			flagID:    "10s",
			status:    http.StatusBadRequest,
			repoError: nil,
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeFlagRepo.repoError = tc.repoError

			w := httptest.NewRecorder()
			req := httptest.NewRequest("DELETE", fmt.Sprintf("/v1/flag/%s", tc.flagID), nil)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code, tc.name)

			if tc.status == http.StatusNoContent {
				suite.Equal(tc.flagID, fmt.Sprintf("%d", suite.fakeFlagRepo.toBeDeleteID))
			}
		})
	}
}

func (suite *FlagHandlerSuite) TestUpdateFlag() {
	cases := []struct {
		name      string
		flagID    string
		req       request.UpdateFlagRequest
		status    int
		repoError error
	}{
		{
			name:   "successfully update flag 1",
			flagID: "10",
			req: request.UpdateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: nil,
			status:    http.StatusOK,
		},
		{
			name:   "failed to update flag 1",
			flagID: "10",
			req: request.UpdateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: model.ErrFlagNotFound,
			status:    http.StatusNotFound,
		},
		{
			name:   "failed to update flag 2",
			flagID: "10",
			req: request.UpdateFlagRequest{
				Flag: request.Flag{
					Tags:        []string{"tag"},
					Description: "description",
					Flag:        "flag",
					Segments: []request.Segment{
						{
							Description: "description",
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
								Key: "on",
							},
						},
					},
				},
			},
			repoError: model.ErrInvalidFlagForUpdate,
			status:    http.StatusBadRequest,
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeFlagRepo.repoError = tc.repoError

			data, err := json.Marshal(tc.req)
			suite.NoError(err)

			w := httptest.NewRecorder()
			url := fmt.Sprintf("/v1/flag/%s", tc.flagID)
			req := httptest.NewRequest("PUT", url, bytes.NewReader(data))

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code, tc.name)

			if tc.status == http.StatusOK {
				suite.Equal(tc.flagID, fmt.Sprintf("%d", suite.fakeFlagRepo.toBeUpdateID))
			}
		})
	}
}

func (suite *FlagHandlerSuite) TestFindByID() {
	cases := []struct {
		name      string
		flagID    string
		status    int
		repoError error
		resp      response.Flag
	}{
		{
			name:      "successfully find flag by its id 1",
			flagID:    "10",
			status:    http.StatusOK,
			repoError: nil,
			resp: response.Flag{
				ID:          10,
				Tags:        []string{"tag1"},
				Description: "description 1",
				Flag:        "flag1",
				Segments: []response.Segment{
					{
						Description: "segment 1",
						Constraints: map[string]response.Constraint{
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
						Variant: response.Variant{
							Key: "on",
						},
					},
				},
			},
		},
		{
			name:      "failed to find flag by its id 1",
			flagID:    "11",
			status:    http.StatusNotFound,
			repoError: nil,
		},
		{
			name:      "failed to find flag by its id 1",
			flagID:    "10",
			status:    http.StatusInternalServerError,
			repoError: errors.New("fake flag repo error"),
		},
	}

	for i := range cases {
		tc := cases[i]
		suite.Run(tc.name, func() {
			suite.fakeFlagRepo.repoError = tc.repoError

			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", fmt.Sprintf("/v1/flag/%s", tc.flagID), nil)

			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			suite.engine.ServeHTTP(w, req)
			suite.Equal(tc.status, w.Code, tc.name)

			if tc.status == http.StatusOK {
				var resp response.Flag

				suite.NoError(json.Unmarshal(w.Body.Bytes(), &resp))

				suite.Equal(tc.resp.ID, resp.ID)
				suite.Equal(tc.resp.Description, resp.Description)
				suite.Equal(tc.resp.Flag, resp.Flag)
				suite.Equal(tc.resp.Tags, resp.Tags)
				suite.Equal(tc.resp.Segments[0].Description, resp.Segments[0].Description)
				suite.Equal(tc.resp.Segments[0].Expression, resp.Segments[0].Expression)
				suite.Equal(tc.resp.Segments[0].Variant.Key, resp.Segments[0].Variant.Key)
			}
		})
	}
}

func TestFlagHandlerSuite(t *testing.T) {
	suite.Run(t, new(FlagHandlerSuite))
}
