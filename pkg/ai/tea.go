package ai

import (
	// "encoding/xml"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	// "github.com/charmbracelet/lipgloss"
	"github.com/alecthomas/chroma/quick"
)

const defaultPrompt = `<role>
You are a software engineer with 200 years experience in writing self-documenting, concise code. You prefer composable code and hold to the tenets of pragramatic programming.
</role>

<instructions>
When providing a snippet of code, always use a code block with the name of the file as a comment at the top. In general, when providing a modification of existing files or previously provided code only provide the diff of the code and if necessary, provide enough of the surrounding code to make the code snippet's intended location clear.

When you generate the answer, first think how the output should be structured and add your answer in <thinking></thinking> tags. Make sure and think step-by-step. This is a space for you to write down relevant content and will not be shown to the user. Use it to consider how the code should be architected, potential sources of bugs, and testing strategies. Once you are done thinking, answer the question. Put your answer inside <answer></answer> XML tags.

Complete the task only if you can produce high-quality code; otherwise tell me you can't and what keeps you from doing so.
</instructions>

`

type state int

const (
	prePromptState state = iota
	contextState
	taskState
	sendState
	responseState
)

type item struct {
	path     string
	selected bool
}

func (i item) Title() string       { return i.path }
func (i item) Description() string { return "" }
func (i item) FilterValue() string { return i.path }

type AIModel struct {
	state       state
	prePrompt   textinput.Model
	contextList list.Model
	task        textinput.Model
	prompt      string
	response    string
	err         error
}

func NewAIModel() AIModel {
	ti := textinput.New()
	ti.Placeholder = "Enter pre-prompt or file path (leave empty for default)"
	ti.Focus()

	items := []list.Item{}
	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Select files/directories for context (space to select, enter when done)"

	task := textinput.New()
	task.Placeholder = "Enter task for Claude"

	return AIModel{
		state:       prePromptState,
		prePrompt:   ti,
		contextList: l,
		task:        task,
	}
}

func (m AIModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			return m.handleEnter()
		case " ":
			if m.state == contextState {
				return m.toggleSelectedItem()
			}
		}

	case promptMsg:
		m.prompt = string(msg)
		return m, sendToClaudeAPI(m.prompt)

	case claudeResponse:
		m.response = string(msg)
		m.state = responseState
		return m, nil
	}

	switch m.state {
	case prePromptState:
		m.prePrompt, cmd = m.prePrompt.Update(msg)
	case contextState:
		m.contextList, cmd = m.contextList.Update(msg)
	case taskState:
		m.task, cmd = m.task.Update(msg)
	}

	return m, cmd
}

func (m AIModel) handleEnter() (tea.Model, tea.Cmd) {
	switch m.state {
	case prePromptState:
		m.state = contextState
		return m, m.loadContextItems
	case contextState:
		m.state = taskState
		m.task.Focus()
		return m, textinput.Blink
	case taskState:
		m.state = sendState
		return m, m.generatePrompt
	}
	return m, nil
}

func (m AIModel) toggleSelectedItem() (tea.Model, tea.Cmd) {
	item := m.contextList.SelectedItem().(item)
	item.selected = !item.selected
	m.contextList.SetItem(m.contextList.Index(), item)
	return m, nil
}

func (m AIModel) loadContextItems() tea.Msg {
	items := []list.Item{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			items = append(items, item{path: path})
		}
		return nil
	})
	if err != nil {
		return errMsg(err)
	}
	m.contextList.SetItems(items)
	return nil
}

type promptMsg string

func (m AIModel) generatePrompt() tea.Cmd {
	return func() tea.Msg {
		prePrompt := m.prePrompt.Value()
		if prePrompt == "" {
			prePrompt = defaultPrompt
		} else if _, err := os.Stat(prePrompt); err == nil {
			content, err := os.ReadFile(prePrompt)
			if err != nil {
				return errMsg(err)
			}
			prePrompt = string(content)
		}

		var context strings.Builder
		context.WriteString("<context>\n")
		for _, item := range m.contextList.Items() {
			if item.(item).selected {
				content, err := os.ReadFile(item.(item).path)
				if err != nil {
					return errMsg(err)
				}
				context.WriteString(fmt.Sprintf("// %s\n%s\n\n", item.(item).path, string(content)))
			}
		}
		context.WriteString("</context>")

		task := fmt.Sprintf("<task>%s</task>", m.task.Value())

		prompt := fmt.Sprintf("%s\n\n%s\n\n%s", prePrompt, context.String(), task)
		return promptMsg(prompt)
	}
}

type claudeResponse string

func sendToClaudeAPI(prompt string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual API call
		return claudeResponse("This is a mock response from Claude API.")
	}
}

type errMsg error

func (m AIModel) View() string {
	switch m.state {
	case prePromptState:
		return fmt.Sprintf(
			"Enter pre-prompt or file path (leave empty for default):\n\n%s\n\n%s",
			m.prePrompt.View(),
			"(press enter to continue)",
		)
	case contextState:
		return m.contextList.View()
	case taskState:
		return fmt.Sprintf(
			"Enter task for Claude:\n\n%s\n\n%s",
			m.task.View(),
			"(press enter to send)",
		)
	case sendState:
		return "Sending request to Claude API..."
	case responseState:
		var highlighted strings.Builder
		quick.Highlight(&highlighted, m.response, "xml", "terminal", "monokai")
		return fmt.Sprintf("Claude's response:\n\n%s", highlighted.String())
	default:
		return "An error occurred."
	}
}

// func main() {
// 	p := tea.NewProgram(NewAIModel())
// 	if _, err := p.Run(); err != nil {
// 		fmt.Printf("Error: %v", err)
// 		os.Exit(1)
// 	}
// }
