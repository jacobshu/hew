package hash

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewHashCmd() *cobra.Command {
	var chars int
	var hashCmd = &cobra.Command{
		Use:   "hash",
		Short: "hashing utilities",
		Args:  cobra.MatchAll(cobra.ExactArgs(2), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%v", args)
		},
	}

	hashCmd.Flags().IntVarP(&chars, "chars", "c", 4, "number of characters per file")
	return hashCmd
}
