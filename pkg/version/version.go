package version

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

// Following variables are set via -ldflags
//nolint:gochecknoglobals
var (
	// VCSRef represents name of branch at build time
	VCSRef string
	// Version represents git SHA at build time
	Version string
	// Date represents the time of build
	Date string
)

func Validate() error {
	missingFields := []string{}

	for name, value := range map[string]string{"VCSRef": VCSRef, "Version": Version, "Date": Date} {
		if value == "" {
			missingFields = append(missingFields, name)
		}
	}

	if len(missingFields) == 0 {
		return nil
	}

	return errors.New("missing build flags")
}

func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := Validate(); err == nil {
			c.Response().Header().Set("X-OPENFLAG-VCS-REF", VCSRef)
			c.Response().Header().Set("X-OPENFLAG-BUILD-DATE", Date)
			c.Response().Header().Set("X-OPENFLAG-BUILD-VERSION", Version)
		}

		return next(c)
	}
}

func String() string {
	return fmt.Sprintf("VCSRef = %s, BuildVersion = %s, BuildDate = %s", VCSRef, Version, Date)
}
