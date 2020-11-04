package cmd

import (
	"os"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd/migrate"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/cmd/server"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/OpenFlag/OpenFlag/pkg/log"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const ExitFailure = 1

func Execute() {
	cfg := config.Init()

	log.SetupLogger(cfg.Logger.AppLogger)

	var root = &cobra.Command{
		Use:   "openflag",
		Short: "OpenFlag is an open source feature flagging, A/B testing and dynamic configuration service.",
	}

	migrate.Register(root, cfg)
	server.Register(root, cfg)

	if err := root.Execute(); err != nil {
		logrus.Errorf("failed to execute root command: %s", err.Error())
		os.Exit(ExitFailure)
	}
}
