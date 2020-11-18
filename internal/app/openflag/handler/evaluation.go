package handler

import (
	"net/http"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/engine"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/request"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

// EvaluationHandler represents a requests handler for evaluations.
type EvaluationHandler struct {
	Engine     engine.Engine
	EntityRepo model.EntityRepo
}

// Evaluate evaluates some entities using an http request.
// nolint:funlen
func (e EvaluationHandler) Evaluate(c echo.Context) error {
	req := request.EvaluationRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("evaluation handler bind (evaluate): %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("evaluation handler validate (evaluate): %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	entities := []model.Entity{}

	for _, entity := range req.Entities {
		entities = append(entities, model.Entity{
			EntityID:      entity.EntityID,
			EntityType:    entity.EntityType,
			EntityContext: entity.EntityContext,
		})
	}

	if req.UseStoredContexts {
		var err error

		entities, err = e.EntityRepo.Find(entities)
		if err != nil {
			logrus.Errorf("evaluation handler failed to use stored contexts: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	resps := []response.EvaluationResponse{}

	for _, entity := range entities {
		result, err := e.Engine.Evaluate(req.Flags, entity)
		if err != nil {
			logrus.Errorf("evaluation handler failed to evaluate: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		evaluations := []response.Evaluation{}

		for _, evaluation := range result.Evaluations {
			evaluations = append(evaluations, response.Evaluation{
				Flag: evaluation.Flag,
				Variant: response.Variant{
					VariantKey:        evaluation.Variant.VariantKey,
					VariantAttachment: evaluation.Variant.VariantAttachment,
				},
			})
		}

		resps = append(resps, response.EvaluationResponse{
			Entity: response.Entity{
				EntityID:      result.Entity.EntityID,
				EntityType:    result.Entity.EntityType,
				EntityContext: result.Entity.EntityContext,
			},
			Evaluations: evaluations,
		})
	}

	if req.SaveContexts {
		if err := e.EntityRepo.Save(entities); err != nil {
			logrus.Errorf("evaluation handler failed to save contexts: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError)
		}
	}

	return c.JSON(http.StatusOK, resps)
}
