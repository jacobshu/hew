package logread

import (
	"github.com/spf13/cobra"
)

func NewLogReadCmd() *cobra.Command {
	var logreadCmd = &cobra.Command{
		Use:   "logread",
		Short: "log reading helper",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			LogRead(args[0])
		},
	}

	return logreadCmd
}
