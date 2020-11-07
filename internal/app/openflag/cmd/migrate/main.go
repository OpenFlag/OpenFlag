package migrate

import (
	"fmt"

	"github.com/OpenFlag/OpenFlag/pkg/database"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const (
	flagPath = "path"
)

// Register register migrate command for openflag binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Provides DB migration functionality",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString(flagPath)
			if err != nil {
				return fmt.Errorf("error parsing %s flag: %s", flagPath, err.Error())
			}

			if path == "" {
				return fmt.Errorf("%s flag is required", flagPath)
			}

			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			path, err := cmd.Flags().GetString(flagPath)
			if err != nil {
				return err
			}

			dbCfg := cfg.Database

			if err := database.Migrate(dbCfg.Driver, dbCfg.MasterConnStr, path); err != nil {
				return err
			}

			logrus.Info("migrations ran successfully")

			return nil
		},
	}

	cmd.Flags().StringP(flagPath, "p", "", "migration folder path")

	root.AddCommand(cmd)
}
