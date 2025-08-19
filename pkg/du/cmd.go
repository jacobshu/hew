package du

import (
	"github.com/spf13/cobra"
)

func NewDuCmd() *cobra.Command {
	var recursive bool

	var duCmd = &cobra.Command{
		Use:   "du",
		Short: "disk utility because I miss it",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			path := "."
			if len(args) > 0 {
				path = args[0]
			}

			GetDU(path, recursive)
		},
	}

	duCmd.Flags().BoolVarP(&recursive, "recursive", "r", false, "Recursively calculate directory sizes")

	return duCmd
}
