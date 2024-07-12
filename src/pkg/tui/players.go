package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/ls13g12/hockey-app/src/pkg/db"
)

const (
	firstName = iota
	lastName uint = 1
)

type playerSidebarModel struct {
	db *mongo.Database
	choices []db.Player
	cursor int
}

type playerCreateModel struct {
	createFunc func(db *mongo.Database, player db.Player) error
	inputs  []textinput.Model
	focused int
	err     error
}

func (m playerSidebarModel) Init() tea.Cmd {
	return nil
}

func (m playerCreateModel) Init() tea.Cmd {
	return textinput.Blink
}


func (m playerCreateModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd = make([]tea.Cmd, len(m.inputs))

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEnter:
			newPlayer := db.Player{
				FirstName: m.inputs[firstName].Value(),
				LastName: m.inputs[lastName].Value(),
			}
			db.CreatePlayer(globalDB, newPlayer)
			return m, nil
		case tea.KeyShiftTab, tea.KeyCtrlP:
			m.prevInput()
		case tea.KeyTab, tea.KeyCtrlN:
			m.nextInput()
		}
		for i := range m.inputs {
			m.inputs[i].Blur()
		}
		m.inputs[m.focused].Focus()
	}

	for i := range m.inputs {
		m.inputs[i], cmds[i] = m.inputs[i].Update(msg)
	}
	return m, tea.Batch(cmds...)
}

func (m playerSidebarModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
				screen_one := ScreenOne(m.db)
				return RootScreen(m.db).SwitchScreen(&screen_one)
			case tea.KeyUp:
				if m.cursor > 0 {
					m.cursor--
				}
			case tea.KeyDown:
				if m.cursor < len(m.choices) -1 {
					m.cursor++
				}
			}	
		}
	return m, nil
}

func (m playerSidebarModel) View() string {
	var s string
	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
				cursor = ">" 
		}

		playerName := fmt.Sprintf("%s %s", choice.FirstName, choice.LastName)

		s += fmt.Sprintf("%s %s\n", cursor, playerName)
	}

	return s
}

func (m playerCreateModel) View() string {
	return fmt.Sprintf(` %s  %s
%s  %s
	
%s
`,
		inputStyle.Width(40).Render("First Name"),
		inputStyle.Width(40).Render("Last Name"),
		m.inputs[firstName].View(),
		m.inputs[lastName].View(),
		continueStyle.Render("Continue ->"),
	) + "\n"
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

func generateInputs() []textinput.Model {
	var inputs []textinput.Model = make([]textinput.Model, 3)
	const firstName = iota
	inputs[firstName] = textinput.New()
	inputs[firstName].Focus()
	inputs[firstName].CharLimit = 20
	inputs[firstName].Width = 30
	inputs[firstName].Prompt = ""

	const lastName = 1
	inputs[lastName] = textinput.New()
	inputs[lastName].CharLimit = 20
	inputs[lastName].Width = 30
	inputs[lastName].Prompt = ""
	return inputs
}
