package main

// A simple program demonstrating the text input component from the Bubbles
// component library.

import (
	"fmt"
	"log"
	"strings"

	"github.com/DeepAung/qcal/calculator"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	red      = lipgloss.Color("#FF0000")
	hotPink  = lipgloss.Color("#FF06B7")
	darkPink = lipgloss.Color("#79305a")
	darkGray = lipgloss.Color("#767676")
)

func main() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

type (
	errMsg error
)

type model struct {
	textInput  textinput.Model
	calculator *calculator.Calculator
	history    []history
	historyIdx int
	err        error
}

type history struct {
	input   string
	output  string
	isError bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = setColor(">> ", hotPink)
	ti.Focus()

	return model{
		textInput:  ti,
		calculator: calculator.NewCalculator(),
		history:    []history{},
		historyIdx: -1,
		err:        nil,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {

		case tea.KeyCtrlC, tea.KeyCtrlD, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyUp:
			m.historyIdx = max(0, m.historyIdx-1)
			if m.historyIdx == len(m.history) {
				m.textInput.SetValue("")
			} else {
				m.textInput.SetValue(m.history[m.historyIdx].input)
			}
		case tea.KeyDown:
			m.historyIdx = min(m.historyIdx+1, len(m.history))
			if m.historyIdx == len(m.history) {
				m.textInput.SetValue("")
			} else {
				m.textInput.SetValue(m.history[m.historyIdx].input)
			}

		case tea.KeyEnter:
			input := m.textInput.Value()
			var output string
			var isError bool

			result, err := m.calculator.Calculate(input)

			if err != nil {
				output = err.Error()
				isError = true
			} else if result == nil {
				output = ""
				isError = false
			} else {
				output = result.Inspect()
				isError = false
			}

			m.history = append(m.history, history{input: input, output: output, isError: isError})
			m.historyIdx = len(m.history)
			m.textInput.Reset()
			return m, cmd
		}

	case errMsg:
		m.err = msg
		return m, nil
	}

	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) View() string {
	if m.err != nil {
		return fmt.Sprintf("ERROR: %v\n", m.err)
	}

	var historyStr []string
	for _, h := range m.history {
		str := ">> " + h.input
		if h.output != "" {
			var coloredOutput string
			if h.isError {
				coloredOutput = setColor(h.output, red)
			} else {
				coloredOutput = setColor(h.output, darkGray)
			}
			str += "\n" + coloredOutput
		}

		historyStr = append(historyStr, str)
	}

	historyRender := strings.Join(historyStr, "\n")
	if historyRender != "" {
		historyRender += "\n"
	}

	return "Welcome to qcal. Enter math expression. (esc to quit)\n" +
		historyRender +
		m.textInput.View()
}

func setColor(str string, color lipgloss.TerminalColor) string {
	return lipgloss.NewStyle().UnsetForeground().Foreground(color).Render(str)
}
