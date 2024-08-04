package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/ls13g12/hockey-app/src/pkg/db"
)

var (
	inputStyle    = lipgloss.NewStyle().Foreground(hotPink)
)

type CreateModel interface {
	Init() tea.Cmd
	View() string
	Update(msg tea.Msg) (CreateModel, tea.Cmd)
}

type playerCreateModel struct {
	confirmCreate bool
	notificationMessage string
	inputs        []textinput.Model
	focused       int
}

func NewCreateModel(resource string) (tea.Model, error) {
	switch resource {
	case "players":
		playerInputs, err := generatePlayerInputs()
		if err != nil {
			return nil, fmt.Errorf("error fetching players")
		}
		return playerCreateModel{
			inputs:  playerInputs,
			focused: 0,
		}, nil
	default:
		return nil, fmt.Errorf("model for create action not implemented on resource %s", resource)
	}
}

func (m playerCreateModel) Init() tea.Cmd {
	return textinput.Blink
}

func (m playerCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	if m.confirmCreate {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.confirmCreate = false
			case tea.KeyEnter:
				player := db.Player{
					PlayerID: uuid.NewString(),
					FirstName: strings.TrimSpace(m.inputs[firstName].View()),
					LastName: strings.TrimSpace(m.inputs[lastName].View()),
					Created: time.Now(),
				}
				err := db.CreatePlayer(globalDB, player)
				if err != nil {
					m.notificationMessage = "player could not be created"
					return m, nil
				}
				homeModel := NewHomeModel("player successfully created")
				return RootScreen().SwitchScreen(homeModel)
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			homeModel := NewHomeModel("")
			return RootScreen().SwitchScreen(homeModel)
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		case tea.KeyEnter:
			if m.focused == len(m.inputs) - 1 {
				m.confirmCreate = true
			} else {
				m.nextInput()
			}
		}

		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		if !m.confirmCreate {
			m.inputs[m.focused].Focus()
		}
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}

	return m, tea.Batch(cmds...)
}

func (m *playerCreateModel) nextInput() {
	m.focused = (m.focused + 1) % len(m.inputs)
}

func (m *playerCreateModel) prevInput() {
	m.focused--

	if m.focused < 0 {
		m.focused = len(m.inputs) - 1
	}
}

func (m playerCreateModel) View() string {
	inputPlayerView := fmt.Sprintf(
`%s %s
%s%s
%s
%s
%s %s
%s%s
%s
%s
`,
		inputStyle.Width(20).Render("First Name"),
		inputStyle.Width(30).Render("Last Name"),
		m.inputs[firstName].View(),
		m.inputs[lastName].View(),
		inputStyle.Render("Nickname"),
		m.inputs[nickname].View(),
		inputStyle.Width(20).Render("Home Number"),
		inputStyle.Width(20).Render("Away Number"),
		m.inputs[homeShirtNumber].View(),
		m.inputs[awayShirtNumber].View(),
		inputStyle.Render("Date Of Birth (dd/mm/yyyy)"),
		m.inputs[dateOfBirthString].View(),
	)

	var confirmMessage string
	if m.confirmCreate {
		confirmMessage = "Click Enter again to confirm"
		return lipgloss.JoinVertical(lipgloss.Top, notificationBox.Render(confirmMessage), inputPlayerView)
	} else if m.notificationMessage != "" {
		return lipgloss.JoinVertical(lipgloss.Top, notificationBox.Render(m.notificationMessage), inputPlayerView)
	}

	return inputPlayerView
}
