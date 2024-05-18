package cmd

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"

  "hew.jacobshu.dev/pkg/cht"
  "hew.jacobshu.dev/pkg/load"
  "hew.jacobshu.dev/pkg/stardew"
)

func BuildCmdTree() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "hew",
		Short: "A handy haversack with tools ready to hand",
		Args:  cobra.NoArgs,
		RunE:  hewRoot,
	}

	var chtCmd = &cobra.Command{
		Use:   "cht",
		Short: "Get help from cht.sh",
		Args:  cobra.NoArgs,
		Run:   chtRoot,
	}
	rootCmd.AddCommand(chtCmd)

	var loadCmd = &cobra.Command{
		Use:   "load",
		Short: "Configure your system for maximum awesomeness",
		Args:  cobra.NoArgs,
		Run:   loadRoot,
	}
	rootCmd.AddCommand(loadCmd)

  var stardewCmd = &cobra.Command{
		Use:   "stardew",
		Short: "game utils",
		Args:  cobra.NoArgs,
		Run:   stardewRoot,
	}
	rootCmd.AddCommand(stardewCmd)
	return rootCmd
}

func hewRoot(cmd *cobra.Command, args []string) error {
	return cmd.Help()
}

func chtRoot(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(
		cht.InitialChtModel(),
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)
	if _, err := p.Run(); err != nil {
		fmt.Printf("Uh oh, there was an error: %v\n", err)
		os.Exit(1)
	}
}

func loadRoot(cmd *cobra.Command, args []string) {
	p := tea.NewProgram(load.NewLoadModel())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running load:", err)
		os.Exit(1)
	}
}

func stardewRoot(cmd *cobra.Command, args []string) {
  p := tea.NewProgram(stardew.NewStardewModel())
  if _, err := p.Run(); err != nil {
    fmt.Printf("Error running stardew: \n%#v", err)
    os.Exit(1)
  }
}

