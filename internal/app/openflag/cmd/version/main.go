package version

import (
	"github.com/OpenFlag/OpenFlag/pkg/version"
	"github.com/spf13/cobra"
)

// Register registers version command for openflag binary.
func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Print the version of OpenFlag",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println(version.String())
			},
		},
	)
}
