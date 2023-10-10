package main

import (
	"fmt"
	"log"
  "math/rand"
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
  IsFile bool
}

type symlinkConfig struct {
  Version   string
  Dotfiles  []symlink
}

var (
	spinnerStyle  = lipgloss.NewStyle().Foreground(forestfox["cyan"])
	helpStyle     = lipgloss.NewStyle().Foreground(forestfox["green"]).Margin(1, 0)
	dotStyle      = helpStyle.Copy().UnsetMargins()
	durationStyle = dotStyle.Copy()
	appStyle      = lipgloss.NewStyle().Margin(1, 2, 0, 2)
)

type symlinkMsg struct {
	duration time.Duration
	source   string
	target   string
	err      error
}

func (s symlinkMsg) String() string {
	if s.duration == 0 {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("🔗 Linked %s to %s in %s", s.source, s.target,
		durationStyle.Render(s.duration.String()))
}

type loadModel struct {
	spinner  spinner.Model
	symlinksToCreate []symlinkMsg
  symlinksCreated  int 
	quitting bool
}

func newLoadModel(symlinksToCreate []symlinkMsg) loadModel {
	s := spinner.New()
	s.Style = spinnerStyle
  s.Spinner = spinner.Points
	return loadModel{
		spinner:  s,
    symlinksToCreate: symlinksToCreate,
    symlinksCreated: 0,
	}
}

func (m loadModel) Init() tea.Cmd {
	readSymlinkConfig()
  log.Printf("%+v", m.symlinksToCreate) 
	return m.spinner.Tick
}

func (m loadModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		m.quitting = true
		return m, tea.Quit
	case symlinkMsg:
    m.symlinksCreated += 1
		m.symlinksToCreate = append(m.symlinksToCreate[1:], msg)
    if m.symlinksCreated == len(m.symlinksToCreate) {
      return m, tea.Quit
    }
		return m, m.createSymlink
	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
    
		return m, cmd
	default:
    if m.symlinksCreated == 0 {
      return m, m.createSymlink
    }
		return m, nil
	}
}

func (m loadModel) View() string {
	var s string

	if m.quitting {
		s += "Symlinked and loaded"
	} else {
		s += m.spinner.View() + " Linking..."
	}

	s += "\n\n"

	for _, res := range m.symlinksToCreate {
    //log.Printf("%+v => %+v", res.source, res.target)
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

func readSymlinkConfig() []symlinkMsg {
	var conf symlinkConfig
  _, err := toml.Decode(symlinksToml, &conf)

  var s []symlinkMsg
  for _, l := range conf.Dotfiles {
    log.Printf("dotfile: %+v", l)
    n := symlinkMsg{source: l.Source, target: l.Target}
    s = append(s, n) 
  }

  if err != nil {
		log.Printf("error reading toml: %+v", err)
	}

  return s
  //log.Printf("read: %+v", s)
}
func (m *loadModel) createSymlink() tea.Msg {
  pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosecA
  time.Sleep(pause)
  msg := m.symlinksToCreate[0]
  m.symlinksToCreate = m.symlinksToCreate[1:]
  msg.duration = pause
  log.Printf("linking: %+v to %+v in %+v,  total: %+v", msg.source, msg.target, msg.duration, m.symlinksCreated)
  return msg
}

func (m *loadModel) nextSymlink() tea.Msg {
  return m.symlinksToCreate[0]
}
