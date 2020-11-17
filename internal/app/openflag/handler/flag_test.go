package handler_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

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

func TestFlagHandlerSuite(t *testing.T) {
	suite.Run(t, new(FlagHandlerSuite))
}
