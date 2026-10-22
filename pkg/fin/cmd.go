package fin

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
	"hew.jacobshu.dev/pkg/shared"
)

func NewFinCmd() *cobra.Command {
	var finCmd = &cobra.Command{
		Use:   "fin",
		Short: "tools for working with hledger",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	finCmd.AddCommand(NewImportCmd())
	finCmd.AddCommand(NewProcCmd())
	return finCmd
}

func NewImportCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "import <file.csv>",
		Short: "read a csv file and print its rows",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			src := args[0]
			target := args[1]
			srcLines, err := ReadCSV(src)
			if err != nil {
				return err
			}
			targetLines, err := ReadCSV(target)
			if err != nil {
				return err
			}

			diff := shared.Difference(srcLines, targetLines)
			log.Printf("src: %d, target: %d, diff: %d\n", len(srcLines), len(targetLines), len(diff))

			b := []byte(strings.Join(append(srcLines, diff...), "\n"))
			err = os.WriteFile(target, b, 0644)
			if err != nil {
				return fmt.Errorf("Error writing to target: %s: %v", target, err)
			}

			return nil
		},
	}
}

func NewProcCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "proc <year>",
		Short: "process accounts",
		Args:  cobra.RangeArgs(0, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			sh := exec.Command("hledger", "--version")
			sh.Stdout = os.Stdout
			sh.Stderr = os.Stderr
			if err := sh.Run(); err != nil {
				return fmt.Errorf("error running hledger: %w", err)
			}
			return nil
		},
	}
}
