package migrate

import (
	"fmt"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/OpenFlag/OpenFlag/pkg/postgres"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file" // Imported for its side effects

	psql "github.com/golang-migrate/migrate/v4/database/postgres"
)

const (
	flagPath = "path"
)

func main(path string, cfg postgres.Config) error {
	pgDb := postgres.WithRetry(postgres.Create, cfg)

	defer func() {
		if err := pgDb.Close(); err != nil {
			logrus.Errorf("postgres connection close error: %s", err.Error())
		}
	}()

	driver, err := psql.WithInstance(pgDb.DB(), &psql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+path, cfg.DBName, driver)
	if err != nil {
		return err
	}

	if err := m.Up(); err == migrate.ErrNoChange {
		logrus.Info("no change detected. All migrations have already been applied!")
		return nil
	} else if err != nil {
		return err
	}

	return nil
}

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

			if err := main(path, cfg.Postgres.Master); err != nil {
				return err
			}

			cmd.Println("migrations ran successfully")

			return nil
		},
	}

	cmd.Flags().StringP(flagPath, "p", "", "migration folder path")

	root.AddCommand(cmd)
}
