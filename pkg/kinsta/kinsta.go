package kinsta

import (
  "log"
  // "fmt"
  // "http"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"hew.jacobshu.dev/pkg/forestfox"

)

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(forestfox.Theme["cyan"])
	helpStyle     = lipgloss.NewStyle().Foreground(forestfox.Theme["green"]).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type kinstaModel struct {
	spinner          spinner.Model
	quitting         bool
}

func NewKinstaModel() kinstaModel {
	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Points

	return kinstaModel{
		spinner:          s,
	}
}

func (m kinstaModel) Init() tea.Cmd {
  log.Println("kinstaModel Init")
	return m.spinner.Tick
}

func (m kinstaModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
    s := msg.String()
    if s == "q" || s == "esc" || s == "ctrl+c" {
      m.quitting = true
      return m, tea.Quit
    }
    return m, nil 
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)

		return m, cmd
	default:
		return m, nil
	}
}

func (m kinstaModel) View() string {
	var s string

	if m.quitting {
		s += "quitting..."
	} else {
		s += m.spinner.View() + " thinking..."
	}

	s += "\n\n"

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}
