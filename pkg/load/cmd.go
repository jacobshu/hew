package load

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewLoadCmd() *cobra.Command {
	var loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Configure your system for maximum awesomeness",
		Args:  cobra.NoArgs,
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(NewLoadModel())
			if _, err := p.Run(); err != nil {
				fmt.Println("Error running load:", err)
				os.Exit(1)
			}
		},
	}

	return loadCmd
}
