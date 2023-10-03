package main

import (
	"fmt"
  "io/ioutil"
  "log"
	"net/http"
  "strings"

  "github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const (
	language = iota
	query
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(forestfox["magenta"])
	continueStyle = lipgloss.NewStyle().Foreground(forestfox["black"])
  continueFocusStyle = lipgloss.NewStyle().
    //Padding(1).
    Foreground(forestfox["cyan"]).
    Background(forestfox["brightBlack"])
)

type chtModel struct {
	inputs   []textinput.Model
  focused  int
  err      error
	response string
}

type (
  chtshMsg string
)

func (c chtshMsg) String() string {
    return fmt.Sprintf("%s", c)
}

func initialChtModel() chtModel {
  var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[language] = textinput.New()
	inputs[language].Placeholder = "go "
	inputs[language].Focus()
	inputs[language].CharLimit = 44
	inputs[language].Width = 50
	inputs[language].Prompt = ""
	//inputs[language].Validate = ccnValidator

	inputs[query] = textinput.New()
	inputs[query].Placeholder = "slices  "
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
    if m.focused <= len(m.inputs) -1  {
      m.inputs[m.focused].Focus()
    }
  case chtshMsg:
    m.response = msg.String()
    log.Printf("msg: %s, m: %+v", msg.String(), m)
    return m, nil
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
	_, err := glamour.Render("# testing \n >glamour\n\n## Headlines", "dark")
  if err != nil {
    return fmt.Sprintf("error rendering with glamour: %+v", err)
  }

  var buttonStyle lipgloss.Style
  if m.focused == len(m.inputs) {
    buttonStyle = continueFocusStyle
  } else {
    buttonStyle = continueStyle
  }

  return fmt.Sprintf(`
 %s
 %s

 %s
 %s

 %s
`,
		inputStyle.Width(30).Render("Language"),
		m.inputs[language].View(),
		inputStyle.Width(10).Render("Query"),
		m.inputs[query].View(),
		buttonStyle.Render("Continue ->"),
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
		m.focused = len(m.inputs)// - 1
	}
}

func getChtsh(language string, query string) tea.Cmd {
	return func() tea.Msg {
	 q := strings.Replace(query, " ", "+", -1)
   req := fmt.Sprintf("https://cht.sh/%s/%s", language, q)
   log.Printf("req: %+v", req)
   resp, err := http.Get(req)
   if err != nil {
      log.Printf("%+v", err)
   }
   
   body, err := ioutil.ReadAll(resp.Body)
   if err != nil {
      log.Printf("%+v", err)
   }
   
   //defer resp.Body.Close()
   return chtshMsg(string(body))
  }
}


  
