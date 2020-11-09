package migrate

import (
	"errors"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/migrations/bindata/mysql"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/migrations/bindata/postgres"

	bindata "github.com/golang-migrate/migrate/v4/source/go_bindata"

	"github.com/OpenFlag/OpenFlag/pkg/database"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/config"
	"github.com/spf13/cobra"
)

// ErrInvalidDatabaseDriver represents an error for returning when we don't support the database driver.
var ErrInvalidDatabaseDriver = errors.New("invalid database driver")

func main(cfg config.Database) error {
	var source *bindata.AssetSource

	switch cfg.Driver {
	case "postgres":
		source = bindata.Resource(postgres.AssetNames(), postgres.Asset)
	case "mysql":
		source = bindata.Resource(mysql.AssetNames(), mysql.Asset)
	default:
		return ErrInvalidDatabaseDriver
	}

	if err := database.Migrate(source, cfg.MasterConnStr); err != nil {
		return err
	}

	return nil
}

// Register register migrate command for openflag binary.
func Register(root *cobra.Command, cfg config.Config) {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		Run: func(cmd *cobra.Command, args []string) {
			if err := main(cfg.Database); err != nil {
				cmd.PrintErrf("failed to run database migrations: %s\n", err.Error())
				return
			}

			cmd.Println("migrations ran successfully")
		},
	}

	root.AddCommand(cmd)
}
