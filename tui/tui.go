package tui

import (
	"context"
	"fmt"

	openai "gollm/openai"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/log"
	oa "github.com/sashabaranov/go-openai"
)

var (
	inputStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("cyan")).Border(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("pink"))
	outputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("green"))
	errorStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("red"))
)

type model struct {
	input    textinput.Model
	output   string
	viewport viewport.Model
	err      error
	apiKey   string
	client   *openai.Client
}

func initialModel(apiKey string) model {
	input := textinput.New()
	input.Placeholder = "Enter a prompt..."
	input.Focus()

	output := ""

	viewport := viewport.New(50, 20)
	viewport.SetContent(output)

	return model{
		input:    input,
		output:   output,
		viewport: viewport,
		err:      nil,
		apiKey:   apiKey,
		client:   openai.NewClient(apiKey),
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			prompt := m.input.Value()
			m.input.SetValue("")

			completion, err := m.client.CreateChatCompletion(context.Background(), oa.ChatCompletionRequest{
				Model: oa.GPT4,
				Messages: []oa.ChatCompletionMessage{
					{
						Role:    oa.ChatMessageRoleUser,
						Content: prompt,
					},
				},
			})

			if err != nil {
				m.err = err
				return m, nil
			}

			m.output += "\nUser: " + prompt + "\nAssistant: " + completion + "\n"
			m.viewport.SetContent(m.output)
			m.viewport.GotoBottom()
		}

	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width - 4
		m.viewport.Height = msg.Height - inputStyle.GetHeight() - 2
	}

	m.input, cmd = m.input.Update(msg)
	m.viewport, cmd = m.viewport.Update(msg)
	return m, cmd
}

func (m model) View() string {
	var output string

	if m.err != nil {
		output = errorStyle.Render("Error: " + m.err.Error())
	} else {
		output = outputStyle.Render(m.viewport.View())
	}

	return fmt.Sprintf(
		"%s\n%s",
		output,
		inputStyle.Render(m.input.View()),
	) + "\n"
}

func RunTUI(apiKey string) {
	p := tea.NewProgram(initialModel(apiKey))
	if _, err := p.Run(); err != nil {
		log.Fatal("error initalizing TUI", "error", err)
	}
}