package kinsta

import (
	"fmt"
	"log"
	// "fmt"
	// "http"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"hew.jacobshu.dev/pkg/forestfox"
)

const (
  columnKeyName = "name"
  columnKeyStatus = "status"
)

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(forestfox.Theme["cyan"])
	helpStyle     = lipgloss.NewStyle().Foreground(forestfox.Theme["green"]).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type kinstaModel struct {
	spinner  spinner.Model
	quitting bool
	sites    []Site
  table    table.Model
}

func NewKinstaModel() kinstaModel {
	s := spinner.New()
	s.Style = spinnerStyle
	s.Spinner = spinner.Points

	return kinstaModel{
		spinner: s,
	}
}

func (m kinstaModel) Init() tea.Cmd {
	log.Println("kinstaModel Init")
  sites, err := GetSites("fbd13128-664b-4cd3-9f1e-725a1a4d6f54")
  if err != nil {
    log.Fatalf("error in GetSites:\n%#v\n", err)
  }

  columns := []table.Column{
		table.NewColumn(columnKeyName, "Name", 20),
		table.NewColumn(columnKeyStatus, "Status", 10),
  }

  rows := []table.Row{}
  for _, site := range sites {
    fmt.Printf("site:\n%v\n", site)
    rows = append(rows, table.Row{site.DisplayName, site.Status})
  }

  t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		// table.WithHeight(7),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

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
