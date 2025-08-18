package cmd

import (
	"github.com/spf13/cobra"

	"hew.jacobshu.dev/pkg/cht"
	"hew.jacobshu.dev/pkg/hash"
	"hew.jacobshu.dev/pkg/load"
  "hew.jacobshu.dev/pkg/ai"
	"hew.jacobshu.dev/pkg/logread"
)

func BuildCmdTree() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "hew",
		Short: "A handy haversack with tools ready to hand",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.AddCommand(cht.NewChtCmd())
	rootCmd.AddCommand(load.NewLoadCmd())
	rootCmd.AddCommand(hash.NewHashCmd())
	rootCmd.AddCommand(ai.NewAICmd())
	rootCmd.AddCommand(logread.NewLogReadCmd())

	return rootCmd
}





