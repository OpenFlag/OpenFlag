package openflag_test

import (
	"testing"
	"time"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

const (
	waitingTimeForDoingSomeJobs = 10 * time.Second
)

type OpenFlagSuite struct {
	suite.Suite
}

func (suite *OpenFlagSuite) TestScenario() {
	root := cmd.NewRootCommand()
	if root == nil {
		suite.Fail("root command is nil!")
		return
	}

	root.SetArgs([]string{"server"})

	go func() {
		if err := root.Execute(); err != nil {
			suite.Fail("failed to execute root command: %s", err.Error())
		}
	}()

	// Wait for starting the server
	time.Sleep(waitingTimeForDoingSomeJobs)

	logrus.Info("start running integration tests!")
}

func TestOpenFlagSuite(t *testing.T) {
	suite.Run(t, new(OpenFlagSuite))
}
