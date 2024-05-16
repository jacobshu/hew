package cht

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"hew.jacobshu.dev/pkg/forestfox"
)

const (
	language = iota
	query
	useHighPerformanceRenderer = false
)

var (
	inputStyle         = lipgloss.NewStyle().Foreground(forestfox.Theme["magenta"])
	continueStyle      = lipgloss.NewStyle().Foreground(forestfox.Theme["brightBlack"])
	continueFocusStyle = lipgloss.NewStyle().Foreground(forestfox.Theme["cyan"])
)

type errMsg struct{ err error }

func (e errMsg) Error() string { return e.err.Error() }

type chtModel struct {
	inputs   []textinput.Model
	focused  int
	err      error
	response string
	viewport viewport.Model
	content  string
	ready    bool
}

type (
	chtshMsg string
)

func (c chtshMsg) String() string {
	return fmt.Sprintf("%s", string(c))
}

func InitialChtModel() chtModel {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[language] = textinput.New()
	inputs[language].Placeholder = "go "
	inputs[language].Focus()
	inputs[language].CharLimit = 44
	inputs[language].Width = 50
	inputs[language].Prompt = ""

	inputs[query] = textinput.New()
	inputs[query].Placeholder = "slices  "
	inputs[query].CharLimit = 100
	inputs[query].Width = 50
	inputs[query].Prompt = ""

	return chtModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
		ready:   false,
	}
}

func (m chtModel) Init() tea.Cmd {
	return nil
}

func (m chtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))
	var (
		cmd tea.Cmd
		//cmds  []tea.Cmd
	)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs) {
				return m, getChtsh(m.inputs[0].Value(), m.inputs[1].Value())
			}
			m.nextInput()
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		if m.focused <= len(m.inputs)-1 {
			m.inputs[m.focused].Focus()
		}
	case chtshMsg:
		m.content = msg.String()
		m.viewport, cmd = m.viewport.Update(msg)
		log.Printf("chtshMsg: %+v", len(msg))
		//m.response = msg.String()
		return m, cmd

	case tea.WindowSizeMsg:
		log.Print("window resize message")
		headerHeight := 8
		//verticalMarginHeight := 8

		if !m.ready {
			vp := viewport.New(40, 20)
			vp.Style = lipgloss.NewStyle().
				Width(40).
				Height(20).
				BorderStyle(lipgloss.RoundedBorder()).
				BorderForeground(forestfox.Theme["blue"]).
				PaddingRight(2)

			m.viewport = vp
			m.viewport.YPosition = headerHeight
			m.viewport.HighPerformanceRendering = useHighPerformanceRenderer
			m.viewport.SetContent(m.content)
			m.ready = true

			// This is only necessary for high performance rendering
			m.viewport.YPosition = headerHeight + 1
		}
		if useHighPerformanceRenderer {
			// This is needed for high-performance rendering only.
			cmds = append(cmds, viewport.Sync(m.viewport))
		}
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	// Handle keyboard and mouse events in the viewport
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)

	return m, tea.Batch(cmds...)
}

func (m chtModel) View() string {
	log.Printf("frame size: %+v\nstyle w: %+v\nviewport w: %+v",
		m.viewport.Style.GetHorizontalFrameSize(),
		m.viewport.Style.GetWidth(),
		m.viewport.Width)
	var buttonStyle lipgloss.Style
	if m.focused == len(m.inputs) {
		buttonStyle = continueFocusStyle
	} else {
		buttonStyle = continueStyle
	}

	return fmt.Sprintf(`
 %s %s

 %s %s

 %s
%s`,
		inputStyle.Width(10).Render("Language"),
		m.inputs[language].View(),
		inputStyle.Width(10).Render("Query"),
		m.inputs[query].View(),
		buttonStyle.Render("Continue ->"),
		m.viewport.View(),
	) + "\n"
}

// nextInput focuses the next input field
func (m *chtModel) nextInput() {
	m.focused = (m.focused + 1) % (len(m.inputs) + 1)
}

// prevInput focuses the previous input field
func (m *chtModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		// don't need to subtract one to account for 0-index since we focus the
		// submit line too
		m.focused = len(m.inputs) // - 1
	}
}

func getChtsh(language string, query string) tea.Cmd {
	return func() tea.Msg {
		q := strings.Replace(query, " ", "+", -1)
		url := fmt.Sprintf("https://cht.sh/%s/%s", language, q)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return errMsg{err}
		}

		req.Header.Set("User-Agent", "curl")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			log.Printf("%+v", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Printf("%+v", err)
		}

		return chtshMsg(string(body))
	}
}
