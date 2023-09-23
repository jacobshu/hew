package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type mode int

const (
	view mode = iota
	create
	edit
)

type model struct {
	tasks    []task
	selected task
	mode     mode
	err      error
}

func getAllTasks() tea.Msg {
	result, err := devDb.getTasks()
	if err != nil {
		return errMsg{err}
	}

	return allTasksMsg(result)
}

type allTasksMsg []task

type errMsg struct{ err error }

// For messages that contain errors it's often handy to also implement the
// error interface on the message.
func (e errMsg) Error() string { return e.err.Error() }

func (m model) Init() tea.Cmd {
	return getAllTasks
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	m.mode = 0
	switch msg := msg.(type) {
	case allTasksMsg:
		m.tasks = msg
		return m, nil

	case errMsg:
		m.err = msg
		return m, tea.Quit

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC || msg.Type == tea.KeyEsc {
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("\nWe had some trouble: %v\n\n", m.err)
	}

	s := fmt.Sprintf("Mode: ... ")
	if len(m.tasks) > 0 {
		s += fmt.Sprintf("%+v !", m.tasks[0])
	}
	return "\n" + s + "\n\n"
}
