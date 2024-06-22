package kinsta

import (
  "fmt"

  "github.com/spf13/cobra"
	"hew.jacobshu.dev/pkg/shared"
)

func NewKinstaCmd() *cobra.Command {
	var kinstaCmd = &cobra.Command{
		Use:   "kinsta",
		Short: "Orchestrate Kinsta from the comfort of your terminal",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			// data, err := kinsta.IsOperationFinished("")
			data, err := CreateManualBackup("", "test")
			if err != nil {
				fmt.Printf("error\n%#v\n", err)
			}

			shared.Pprint(data)
		},
	}

	return kinstaCmd
}
