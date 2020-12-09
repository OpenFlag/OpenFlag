package grpc

import (
	"context"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/engine"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/evaluation"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/request"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type evaluationServer struct {
	Engine     engine.Engine
	EntityRepo model.EntityRepo
}

// Evaluate evaluates some entities using a gRPC request.
// nolint:funlen
func (s *evaluationServer) Evaluate(
	c context.Context, r *evaluation.EvaluationRequest,
) (*evaluation.EvaluationResponseList, error) {
	req, err := s.validate(r)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
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

		entities, err = s.EntityRepo.Find(entities)
		if err != nil {
			logrus.Errorf("grpc evaluation handler failed to use stored contexts: %s", err.Error())
			return nil, status.Error(codes.Internal, ErrInternalServerError.Error())
		}
	}

	resps := []*evaluation.EvaluationResponse{}

	for _, entity := range entities {
		result, err := s.Engine.Evaluate(req.Flags, entity)
		if err != nil {
			logrus.Errorf("grpc evaluation handler failed to evaluate: %s", err.Error())
			return nil, status.Error(codes.Internal, ErrInternalServerError.Error())
		}

		evaluations := []*evaluation.EvaluationResponse_Evaluation{}

		for _, ev := range result.Evaluations {
			evaluations = append(evaluations, &evaluation.EvaluationResponse_Evaluation{
				Flag: ev.Flag,
				Variant: &evaluation.EvaluationResponse_Variant{
					VariantKey:        ev.Variant.VariantKey,
					VariantAttachment: ev.Variant.VariantAttachment,
				},
			})
		}

		resps = append(resps, &evaluation.EvaluationResponse{
			Entity: &evaluation.Entity{
				EntityID:      result.Entity.EntityID,
				EntityType:    result.Entity.EntityType,
				EntityContext: result.Entity.EntityContext,
			},
			Evaluations: evaluations,
		})
	}

	if req.SaveContexts {
		if err := s.EntityRepo.Save(entities); err != nil {
			logrus.Errorf("grpc evaluation handler failed to save contexts: %s", err.Error())
			return nil, status.Error(codes.Internal, ErrInternalServerError.Error())
		}
	}

	return &evaluation.EvaluationResponseList{List: resps}, nil
}

func (s *evaluationServer) validate(req *evaluation.EvaluationRequest) (request.EvaluationRequest, error) {
	entities := []request.Entity{}

	for i := range req.Entities {
		entity := request.Entity{}

		if req.Entities[i] != nil {
			entity = request.Entity{
				EntityID:      req.Entities[i].EntityID,
				EntityType:    req.Entities[i].EntityType,
				EntityContext: req.Entities[i].EntityContext,
			}
		}

		entities = append(entities, entity)
	}

	result := request.EvaluationRequest{
		Entities:          entities,
		Flags:             req.Flags,
		SaveContexts:      req.SaveContexts,
		UseStoredContexts: req.UseStoredContexts,
	}

	return result, result.Validate()
}
