package cmd

import (
	"os"

	"github.com/OpenFlag/OpenFlag/pkg/version"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd/migrate"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd/server"
	versionCmd "github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd/version"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/OpenFlag/OpenFlag/pkg/log"
	"github.com/spf13/cobra"
)

const (
	exitFailure = 1
)

// Execute executes the main functionality of openflag binary.
func Execute() {
	var root = &cobra.Command{
		Use:   "openflag",
		Short: "OpenFlag is an open-source feature flagging, A/B testing, and dynamic configuration service.",
	}

	cfg := config.Init()

	log.SetupLogger(cfg.Logger.AppLogger)

	if err := version.Validate(); err != nil {
		root.PrintErrln(err.Error())
		return
	}

	versionCmd.Register(root)
	migrate.Register(root, cfg)
	server.Register(root, cfg)

	if err := root.Execute(); err != nil {
		os.Exit(exitFailure)
	}
}
