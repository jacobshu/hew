package cht

import (
  "fmt"
  "os"

  "github.com/spf13/cobra"
  tea "github.com/charmbracelet/bubbletea"
)

func NewChtCmd() *cobra.Command {
	var chtCmd = &cobra.Command{
		Use:   "cht",
		Short: "Get help from cht.sh",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(
				InitialChtModel(),
				tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
				tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
			)
			if _, err := p.Run(); err != nil {
				fmt.Printf("Uh oh, there was an error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	return chtCmd
}
