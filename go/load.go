package main

import (
 	"fmt"
  "log"
	"math/rand"
	"os"
	"strings"
	"time"

  "github.com/BurntSushi/toml"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss" 
)

type symlinkConfig struct {
	source     string
	target     string
	linkType   string
}

func symlink(src string, target string) {
  return
}

func readConfig() {
  return
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
	src     string
  target  string
  err     error
}

func (s symlinkMsg) String() string {
	if s.duration == 0 {
		return dotStyle.Render(strings.Repeat(".", 30))
	}
	return fmt.Sprintf("🍔 Linked %s to %s %s", s.src, s.target,
		durationStyle.Render(s.duration.String()))
}

type loadModel struct {
	spinner  spinner.Model
	symlinks  []symlinkMsg
	quitting bool
}

func newLoadModel() loadModel {
	const numLastResults = 5
	s := spinner.New()
	s.Style = spinnerStyle
  s.Spinner = spinner.Moon
	return loadModel{
		spinner: s,
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
		s += m.spinner.View() + " Eating food..."
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
 links, err := toml.Decode(tomlData, &conf)
 if err != nil {
   log.Printf("error reading toml: %+v", err)
 }


 
}

func main() {
	p := tea.NewProgram(newLoadModel())

	// Simulate activity
	go func() {
		for {
			pause := time.Duration(rand.Int63n(899)+100) * time.Millisecond // nolint:gosec
			time.Sleep(pause)

			// Send the Bubble Tea program a message from outside the
			// tea.Program. This will block until it is ready to receive
			// messages.
			p.Send(symlinkMsg{src: "", duration: pause})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func randomFood() string {
	food := []string{
		"an apple", "a pear", "a gherkin", "a party gherkin",
		"a kohlrabi", "some spaghetti", "tacos", "a currywurst", "some curry",
		"a sandwich", "some peanut butter", "some cashews", "some ramen",
	}
	return food[rand.Intn(len(food))] // nolint:gosec
}
