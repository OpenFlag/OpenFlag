package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/request"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/response"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
)

var (
	// ErrInvalidJSONSyntax represents an error that we return when we can not parse the given json request.
	ErrInvalidJSONSyntax = errors.New("invalid json syntax")
)

// FlagHandler represents a requests handler for flags.
type FlagHandler struct {
	FlagRepo model.FlagRepo
}

// Create creates a flag using http request.
func (f FlagHandler) Create(c echo.Context) error {
	req := request.CreateFlagRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("flag handler bind data: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("flag handler validation: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	flag, err := f.flagFromRequest(req.Flag)
	if err != nil {
		logrus.Errorf("flag handler flag from request failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := f.FlagRepo.Create(flag); err != nil {
		if err == model.ErrDuplicateFlagFound {
			return echo.NewHTTPError(http.StatusConflict, err.Error())
		}

		logrus.Errorf("flag handler failed (create): %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resp, err := f.responseFromFlag(flag)
	if err != nil {
		logrus.Errorf("flag handler response from flag failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, *resp)
}

func (f FlagHandler) flagFromRequest(req request.Flag) (*model.Flag, error) {
	segments := []model.Segment{}

	for _, segment := range req.Segments {
		constraints := map[string]model.Constraint{}

		for identifier, constraint := range segment.Constraints {
			constraints[identifier] = model.Constraint{
				Name:       constraint.Name,
				Parameters: constraint.Parameters,
			}
		}

		segments = append(segments, model.Segment{
			Description: segment.Description,
			Constraints: constraints,
			Expression:  segment.Expression,
			Variant: model.Variant{
				Key:        segment.Variant.Key,
				Attachment: segment.Variant.Attachment,
			},
		})
	}

	var tags *string = nil

	if len(req.Tags) != 0 {
		tagsBytes, err := json.Marshal(req.Tags)
		if err != nil {
			return nil, err
		}

		tagsStr := string(tagsBytes)
		tags = &tagsStr
	}

	segmentsByte, err := json.Marshal(segments)
	if err != nil {
		return nil, err
	}

	segmentsStr := string(segmentsByte)

	flag := model.Flag{
		Tags:        tags,
		Description: req.Description,
		Flag:        req.Flag,
		Segments:    segmentsStr,
	}

	return &flag, nil
}

func (f FlagHandler) responseFromFlag(flag *model.Flag) (*response.Flag, error) {
	var flagSegments []model.Segment

	if err := json.Unmarshal([]byte(flag.Segments), &flagSegments); err != nil {
		return nil, err
	}

	segments := []response.Segment{}

	for _, segment := range flagSegments {
		constraints := map[string]response.Constraint{}

		for identifier, constraint := range segment.Constraints {
			constraints[identifier] = response.Constraint{
				Name:       constraint.Name,
				Parameters: constraint.Parameters,
			}
		}

		segments = append(segments, response.Segment{
			Description: segment.Description,
			Constraints: constraints,
			Expression:  segment.Expression,
			Variant: response.Variant{
				Key:        segment.Variant.Key,
				Attachment: segment.Variant.Attachment,
			},
		})
	}

	var tags []string = nil

	if flag.Tags != nil {
		if err := json.Unmarshal([]byte(*flag.Tags), &tags); err != nil {
			return nil, err
		}
	}

	resp := response.Flag{
		ID:          flag.ID,
		Tags:        tags,
		Description: flag.Description,
		Flag:        flag.Flag,
		Segments:    segments,
		CreatedAt:   flag.CreatedAt,
		DeletedAt:   flag.DeletedAt,
	}

	return &resp, nil
}
