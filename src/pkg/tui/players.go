package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

type screenTwoModel struct {
	choices []string
}

func ScreenTwo() screenTwoModel {
	return screenTwoModel{
		choices: []string{"player1", "player2"},
	}
}

func (m screenTwoModel) Init() tea.Cmd {
	return nil
}


func (m screenTwoModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC:
					return m, tea.Quit
			default:
				screen_one := ScreenOne()
				return RootScreen().SwitchScreen(&screen_one)
			}
	default:
			return m, nil
	}
}

func (m screenTwoModel) View() string {
	str := "This is the second screen. Press any key to switch to the first screen."
	return str
}
