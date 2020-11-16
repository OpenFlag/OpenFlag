package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

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

// Create creates a flag using an http request.
func (f FlagHandler) Create(c echo.Context) error {
	req := request.CreateFlagRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("flag handler bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("flag handler validate: %s", err.Error())
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

// Delete deletes a flag using an http request.
func (f FlagHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		logrus.Errorf("flag handler param: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	if err := f.FlagRepo.Delete(id); err != nil {
		logrus.Errorf("flag handler failed (delete): %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusNoContent)
}

// Update updates a flag using an http request.
func (f FlagHandler) Update(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		logrus.Errorf("flag handler param: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	req := request.UpdateFlagRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("flag handler bind data: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("flag handler validate: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	flag, err := f.flagFromRequest(req.Flag)
	if err != nil {
		logrus.Errorf("flag handler flag from request failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if err := f.FlagRepo.Update(id, flag); err != nil {
		if err == model.ErrFlagNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		if err == model.ErrInvalidFlagForUpdate {
			return echo.NewHTTPError(http.StatusBadRequest, err.Error())
		}

		logrus.Errorf("flag handler failed (update): %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resp, err := f.responseFromFlag(flag)
	if err != nil {
		logrus.Errorf("flag handler response from flag failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, *resp)
}

// FindByID finds a flag by its given id using an http request.
func (f FlagHandler) FindByID(c echo.Context) error {
	id, err := strconv.ParseInt(c.Param("id"), 0, 64)
	if err != nil {
		logrus.Errorf("flag handler param: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	flag, err := f.FlagRepo.FindByID(id)
	if err != nil {
		if err == model.ErrFlagNotFound {
			return echo.NewHTTPError(http.StatusNotFound, err.Error())
		}

		logrus.Errorf("flag handler failed (find by id): %s", err.Error())

		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resp, err := f.responseFromFlag(flag)
	if err != nil {
		logrus.Errorf("flag handler response from flag failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, *resp)
}

// FindByTag finds flags that hav given tag using an http request.
func (f FlagHandler) FindByTag(c echo.Context) error {
	req := request.FindFlagsByTagRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("flag handler bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("flag handler validate: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	flags, err := f.FlagRepo.FindByTag(req.Tag)
	if err != nil {
		logrus.Errorf("flag handler failed (find by tag): %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resp, err := f.responseFromFlags(flags)
	if err != nil {
		logrus.Errorf("flag handler response from flags failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, resp)
}

// FindByFlag finds history of a flag using an http request.
func (f FlagHandler) FindByFlag(c echo.Context) error {
	req := request.FindFlagHistoryRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("flag handler bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("flag handler validate: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	flags, err := f.FlagRepo.FindByFlag(req.Flag)
	if err != nil {
		logrus.Errorf("flag handler failed (find by flag): %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resp, err := f.responseFromFlags(flags)
	if err != nil {
		logrus.Errorf("flag handler response from flags failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, resp)
}

// FindFlags finds flags with offset and limit using an http request.
func (f FlagHandler) FindFlags(c echo.Context) error {
	req := request.FindFlagsRequest{}

	if err := c.Bind(&req); err != nil {
		logrus.Errorf("flag handler bind: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, ErrInvalidJSONSyntax.Error())
	}

	if err := req.Validate(); err != nil {
		logrus.Errorf("flag handler validate: %s", err.Error())
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	timestamp := time.Now()
	if req.Timestamp != nil && !req.Timestamp.IsZero() {
		timestamp = *req.Timestamp
	}

	flags, err := f.FlagRepo.FindFlags(req.Offset, req.Limit, timestamp)
	if err != nil {
		logrus.Errorf("flag handler failed (find flags): %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	resp, err := f.responseFromFlags(flags)
	if err != nil {
		logrus.Errorf("flag handler response from flags failed: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, resp)
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

func (f FlagHandler) responseFromFlags(flags []model.Flag) ([]response.Flag, error) {
	resps := []response.Flag{}

	for _, flag := range flags {
		fl := flag

		resp, err := f.responseFromFlag(&fl)
		if err != nil {
			return nil, err
		}

		resps = append(resps, *resp)
	}

	return resps, nil
}
