package version

import (
	"github.com/OpenFlag/OpenFlag/pkg/version"
	"github.com/spf13/cobra"
)

func Register(root *cobra.Command) {
	root.AddCommand(
		&cobra.Command{
			Use:   "version",
			Short: "Print the version number of OpenFlag",
			Run: func(cmd *cobra.Command, args []string) {
				cmd.Println(version.String())
			},
		},
	)
}
