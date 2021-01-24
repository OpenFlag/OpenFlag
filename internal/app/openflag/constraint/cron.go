package constraint

import (
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/robfig/cron/v3"
)

const parserFormat = cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow

// CronConstraint represents Openflag cron constraint.
type CronConstraint struct {
	Expression string `json:"expression"`
	schedule   cron.Schedule
}

// Name is an implementation for the Constraint interface.
func (c CronConstraint) Name() string {
	return CronConstraintName
}

// Validate is an implementation for the Constraint interface.
func (c CronConstraint) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Expression, validation.By(func(v interface{}) error {
			parser := cron.NewParser(parserFormat)
			_, err := parser.Parse(v.(string))

			return err
		})),
	)
}

// Initialize is an implementation for the Constraint interface.
func (c *CronConstraint) Initialize() error {
	parser := cron.NewParser(parserFormat)

	schedule, err := parser.Parse(c.Expression)
	if err != nil {
		return err
	}

	c.schedule = schedule

	return nil
}

// Evaluate is an implementation for the Constraint interface.
func (c CronConstraint) Evaluate(_ model.Entity) bool {
	now := time.Now().Truncate(time.Minute)

	next := c.schedule.Next(now.Add(-1 * time.Minute))

	return now.Equal(next)
}
