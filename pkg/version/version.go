package version

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
)

// Following variables are set via -ldflags
// nolint:gochecknoglobals
var (
	// AppVersion represents the semantic version of the app
	AppVersion string
	// VCSRef represents name of branch at build time
	VCSRef string
	// Version represents git SHA at build time
	BuildVersion string
	// Date represents the time of build
	Date string
)

// Validate validates the version variables.
func Validate() error {
	missingFields := []string{}

	for name, value := range map[string]string{
		"AppVersion": AppVersion, "VCSRef": VCSRef, "BuildVersion": BuildVersion, "Date": Date,
	} {
		if value == "" {
			missingFields = append(missingFields, name)
		}
	}

	if len(missingFields) == 0 {
		return nil
	}

	return errors.New("missing build flags")
}

// Middleware is an echo middleware for adding version variables to the response header of OpenFlag.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if err := Validate(); err == nil {
			c.Response().Header().Set("X-OPENFLAG-APP-VERSION", AppVersion)
			c.Response().Header().Set("X-OPENFLAG-VCS-REF", VCSRef)
			c.Response().Header().Set("X-OPENFLAG-BUILD-VERSION", BuildVersion)
			c.Response().Header().Set("X-OPENFLAG-BUILD-DATE", Date)
		}

		return next(c)
	}
}

// String generates a string for representing the version of OpenFlag.
func String() string {
	return fmt.Sprintf(
		"AppVersion = %s\n"+
			"VCSRef = %s\n"+
			"BuildVersion = %s\n"+
			"BuildDate = %s",
		AppVersion, VCSRef, BuildVersion, Date,
	)
}
