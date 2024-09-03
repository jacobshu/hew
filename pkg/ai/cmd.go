package ai

import (
	"fmt"
	"os"
	
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

func NewAICmd() *cobra.Command {
	var model string
	var context string
	var aiCmd = &cobra.Command{
		Use:   "ai",
		Short: "Ask AI if you want ¯\\_(ツ)_/¯",
		Args:  cobra.MatchAll(cobra.ArbitraryArgs, cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			p := tea.NewProgram(NewAIModel())
			if _, err := p.Run(); err != nil {
				fmt.Printf("Alas, there's been an error: %v", err)
				os.Exit(1)
			}

		},
	}

	aiCmd.Flags().StringVarP(&model, "model", "m", "sonnet", "which claude model to use")
	aiCmd.Flags().StringVarP(&context, "context", "c", "./", "a group of files to include as context to the prompt")
	return aiCmd
}
