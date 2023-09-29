package main

import (
	"fmt"
	"net/http"

  "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	//"github.com/charmbracelet/lipgloss"
)

const (
	language = iota
	query
)

type chtModel struct {
	inputs   []textinput.Model
  focused  int
  err      error
	query    string
	request  string
	response http.Response
}

func initialChtModel() chtModel {
  var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[language] = textinput.New()
	inputs[language].Placeholder = "language"
	inputs[language].Focus()
	inputs[language].CharLimit = 44
	inputs[language].Width = 50
	inputs[language].Prompt = ""
	//inputs[language].Validate = ccnValidator

	inputs[query] = textinput.New()
	inputs[query].Placeholder = "query "
	inputs[query].CharLimit = 100
	inputs[query].Width = 50
	inputs[query].Prompt = ""
	//inputs[query].Validate = expValidator

	return chtModel{
		inputs:  inputs,
		focused: 0,
		err:     nil,
	}
}

func (m chtModel) Init() tea.Cmd {
	return nil
}

func (m chtModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			if m.focused == len(m.inputs)-1 {
				return m, tea.Quit
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
		m.inputs[m.focused].Focus()

	// We handle errors just like any other message
	case errMsg:
		m.err = msg
		return m, nil
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m chtModel) View() string {
	out, err := glamour.Render("", "dark")
	return fmt.Sprintf("")
}

// nextInput focuses the next input field
func (m *chtModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

// prevInput focuses the previous input field
func (m *chtModel) prevInput() {
	m.focused--
	// Wrap around
	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}
