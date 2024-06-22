package hash

import (
	"fmt"

	"github.com/spf13/cobra"
)

func NewHashCmd() *cobra.Command {
	var hashCmd = &cobra.Command{
		Use:   "hash",
		Short: "hashing utilities",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Calling hash")
		},
	}

  hashCmd.AddCommand(NewTestCmd())
	return hashCmd
}

func NewTestCmd() *cobra.Command {
	var chars int
	var testCmd = &cobra.Command{
		Use:   "test",
		Short: "testing utilities for the hash command",
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			chars, _ := cmd.Flags().GetInt("chars")
			fmt.Printf("%v", chars)
		},
	}

	testCmd.Flags().IntVarP(&chars, "chars", "c", 4, "number of characters per file")
	return testCmd
}
