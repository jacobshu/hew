package main

import (
	"fmt"
	"log"

	//"os"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type symlink struct {
  Source string 
  Target string
  IsFile string
}

type symlinkConfig struct {
  Version   string
  Dotfiles  []symlink
}

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(forestfox["cyan"])
	helpStyle     = lipgloss.NewStyle().Foreground(forestfox["black"]).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type symlinkMsg struct {
	duration time.Duration
	src      string
	target   string
	err      error
}

func (s symlinkMsg) String() string {
	if s.duration == 0 {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("💾 Linked %s to %s in %s", s.src, s.target,
		durationStyle.Render(s.duration.String()))
}

type loadModel struct {
	spinner  spinner.Model
	symlinks []symlinkMsg
	quitting bool
}

func newLoadModel() loadModel {
	const numLastResults = 5
	s := spinner.New()
	s.Style = spinnerStyle
  s.Spinner = spinner.Points
	return loadModel{
		spinner:  s,
		symlinks: make([]symlinkMsg, numLastResults),
	}
}

func (m loadModel) Init() tea.Cmd {
	m.readConfig()
	return m.spinner.Tick
}

func (m loadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case symlinkMsg:
		m.symlinks = append(m.symlinks[1:], msg)
		return m, nil
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	default:
		return m, nil
	}
}

func (m loadModel) View() string {
	var s string

	if m.quitting {
		s += "That’s all for today!"
	} else {
		s += m.spinner.View() + " Linking..."
	}

	s += "\n\n"

	for _, res := range m.symlinks {
		s += res.String() + "\n"
	}

	if !m.quitting {
		s += helpStyle.Render("Press any key to exit")
	}

	if m.quitting {
		s += "\n"
	}

	return appStyle.Render(s)
}

func (*loadModel) readConfig() {
	var conf symlinkConfig
	md, err := toml.Decode(symlinksToml, &conf)
  log.Printf("%+v", md.Undecoded())
	log.Printf("links: %+v", conf)

	if err != nil {
		log.Printf("error reading toml: %+v", err)
	}
}

