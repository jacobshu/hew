package cht

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  tea "charm.land/bubbletea/v2"
)

func NewChtCmd() *cobra.Command {
	var chtCmd = &cobra.Command{
		Use:   "cht",
		Short: "Get help from cht.sh",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(
				InitialChtModel(),
			)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Uh oh, there was an error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return chtCmd
}
