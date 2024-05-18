package stardew

import (
	"encoding/xml"
	"fmt"
  "io/fs"
	"log"
	"os"
	"path"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"hew.jacobshu.dev/pkg/forestfox"
)

type stardewModel struct {
	save_dirs []string
	quitting  bool
}

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(forestfox.Theme["cyan"])
	helpStyle     = lipgloss.NewStyle().Foreground(forestfox.Theme["green"]).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

func NewStardewModel() stardewModel {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}

	savesDir := path.Join(homeDir, ".config/StardewValley/Saves")
	files, err := os.ReadDir(savesDir)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		fmt.Println(file)
	}

	fs.WalkDir(os.DirFS(savesDir), ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Fatal(err)
		}

    info, err := d.Info()
    if err != nil {
      log.Printf("error getting info: %#v", err)
    }
    fmt.Printf("path: %#v, \ninfo: %#v\n", path, info)
		return nil
	})

	return stardewModel{}
}

func IsXML(data []byte) bool {
    return xml.Unmarshal(data, new(interface{})) == nil
}

func (m stardewModel) Init() tea.Cmd {
	log.Println("stardewModel Init")
	return nil
}

func (m stardewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		s := msg.String()
		if s == "q" || s == "esc" || s == "ctrl+c" {
			m.quitting = true
			return m, tea.Quit
		}
		return m, nil
	default:
		return m, nil
	}
}

func (m stardewModel) View() string {
	var s string

	if m.quitting {
		s += "Until next time..."
	} else {
		s += " Reading..."
	}

	s += "\n\n"

	if !m.quitting {
		s += helpStyle.Render("Press \"q\" or \"esc\" to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}
