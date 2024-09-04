package ai

import (
	// "encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/sabhiram/go-gitignore"
)

const defaultPrompt = `<role>
You are a software engineer with 2000 years experience in writing self-documenting, concise code. You prefer composable code and hold to the tenets of pragramatic programming.
</role>

<instructions>
When providing a snippet of code, always use a code block with the name of the file as a comment at the top. In general, when providing a modification of existing files or previously provided code only provide the diff of the code and if necessary, provide enough of the surrounding code to make the code snippet's intended location clear.

When you generate the answer, first think how the output should be structured and add your answer in <thinking></thinking> tags. Make sure and think step-by-step. This is a space for you to write down relevant content and will not be shown to the user. Use it to consider how the code should be architected, potential sources of bugs, and testing strategies. Once you are done thinking, answer the question. Put your answer inside <answer></answer> XML tags.

Complete the task only if you can produce high-quality code; otherwise tell me you can't and what keeps you from doing so.
</instructions>`

const (
	prePromptState state = iota
	contextState
	taskState
	sendState
	responseState
)

var (
	itemStyle         = lipgloss.NewStyle().PaddingLeft(4)
	selectedItemStyle = lipgloss.NewStyle().PaddingLeft(2).Foreground(lipgloss.Color("170"))
)

type claudeResponse string
type errMsg error
type itemDelegate struct {
  list.DefaultDelegate
}
type loadedItemsMsg struct {
  items []list.Item
}
type promptMsg string
type state int

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

	l := list.New([]list.Item{}, newItemDelegate(), 0, 0)
	l.Title = "Select files/directories for context (space to select, enter when done)"
	l.SetShowStatusBar(false)
	l.SetFilteringEnabled(false)
	// l.Styles.PaginationStyle = list.Pagination{}.Style().PaddingLeft(4)
	l.Styles.HelpStyle = list.DefaultStyles().HelpStyle.PaddingLeft(4).PaddingBottom(1)

	task := textinput.New()
	task.Placeholder = "Enter task for Claude"

	return AIModel{
		state:       prePromptState,
		prePrompt:   ti,
		contextList: l,
		task:        task,
	}
}

func newItemDelegate() *itemDelegate {
	d := &itemDelegate{}

	d.UpdateFunc = func(msg tea.Msg, m *list.Model) tea.Cmd {
		if msg, ok := msg.(tea.KeyMsg); ok && msg.String() == " " {
			index := m.Index()
			if item, ok := m.SelectedItem().(item); ok {
				item.selected = !item.selected
				m.SetItem(index, item)
			}
		}
		return nil
	}

	return d
}

func (d itemDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	i, ok := listItem.(item)
	if !ok {
		return
	}

	str := i.Title()
	var renderedStr string
	if index == m.Index() {
		renderedStr = selectedItemStyle.Render("> " + str)
	} else if i.selected {
		renderedStr = selectedItemStyle.Render("[x] " + str)
	} else {
		renderedStr = itemStyle.Render(str)
	}
	fmt.Fprint(w, renderedStr)
}

func (m AIModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m AIModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case loadedItemsMsg:
		m.contextList.SetItems(msg.items)
		return m, nil
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
		return m, tea.Batch(m.loadContextItems(), m.contextList.StartSpinner())
	case contextState:
		m.state = taskState
		m.task.Focus()
		return m, textinput.Blink
	case taskState:
		m.state = sendState
		return m, m.generatePrompt()
	}
	return m, nil
}

func (m AIModel) toggleSelectedItem() (tea.Model, tea.Cmd) {
	if i, ok := m.contextList.SelectedItem().(item); ok {
		i.selected = !i.selected
		m.contextList.SetItem(m.contextList.Index(), i)
	}
	return m, nil
}

func (m *AIModel) loadContextItems() tea.Cmd {
	return func() tea.Msg {
		items := []list.Item{}

		// Load .gitignore if it exists
		ignore, _ := ignore.CompileIgnoreFile(".gitignore")

		err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

      if strings.HasPrefix(path, ".git") {
        return nil
      }

			// Check if the file should be ignored
			if ignore != nil && ignore.MatchesPath(path) {
				if info.IsDir() {
					return filepath.SkipDir
				}
				return nil
			}

			items = append(items, item{path: path, selected: false})
			return nil
		})
		if err != nil {
			return errMsg(err)
		}
		return loadedItemsMsg{items: items}
	}
}

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
		for _, listItem := range m.contextList.Items() {
			if i, ok := listItem.(item); ok && i.selected {
				content, err := os.ReadFile(i.path)
				if err != nil {
					return errMsg(err)
				}
				context.WriteString(fmt.Sprintf("// %s\n%s\n\n", i.path, string(content)))
			}
		}
		context.WriteString("</context>")

		task := fmt.Sprintf("<task>%s</task>", m.task.Value())

		prompt := fmt.Sprintf("%s\n\n%s\n\n%s", prePrompt, context.String(), task)
		return promptMsg(prompt)
	}
}


func sendToClaudeAPI(prompt string) tea.Cmd {
	return func() tea.Msg {
		// TODO: Implement actual API call
		return claudeResponse(fmt.Sprintf("prompt: %v\n", prompt))
	}
}


func (m AIModel) View() string {
	switch m.state {
	case prePromptState:
		return fmt.Sprintf(
			"Enter pre-prompt or file path (leave empty for default):\n\n%s\n\n%s",
			m.prePrompt.View(),
			"(press enter to continue)",
		)
	case contextState:
		if len(m.contextList.Items()) == 0 {
			return "Loading items..."
		}
		return fmt.Sprintf(
			"Select files/directories for context:\n\n%s\n\n(space to select, enter when done)",
			m.contextList.View(),
		)
	case taskState:
		return fmt.Sprintf(
			"Enter task for Claude:\n\n%s\n\n%s",
			m.task.View(),
			"(press enter to send)",
		)
	case sendState:
    return fmt.Sprintf("model: %v\n", m) //"Sending request to Claude API..."
	case responseState:
		var highlighted strings.Builder
		quick.Highlight(&highlighted, m.response, "xml", "terminal", "monokai")
		return fmt.Sprintf("Claude's response:\n\n%s", highlighted.String())
	default:
		return "An error occurred."
	}
}
