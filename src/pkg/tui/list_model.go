package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ls13g12/hockey-app/src/pkg/db"
)

const (
	hotPink  = lipgloss.Color("#FF06B7")
	darkGray = lipgloss.Color("#767676")
)

var (
	listBox 		 = lipgloss.NewStyle().Height(20).Margin(0, 2, 0, 0).Padding(0, 2, 0, 0).BorderStyle(lipgloss.NormalBorder()).BorderRight(true)
	headingStyle = lipgloss.NewStyle().Foreground(hotPink)
	textStyle    = lipgloss.NewStyle().Foreground(darkGray)
)

type ListModel interface {
	Init() tea.Cmd
	View() string
	Update(msg tea.Msg) (ListModel, tea.Cmd)
}

type playerListModel struct {
	choices []db.Player
	cursor	int
}

func NewListModel(resource string) playerListModel {
	switch resource {
		case "players":
			players, err := db.AllPlayers(globalDB)
			if err != nil {
				panic(err)
			}
			return playerListModel{
				choices: players,
				cursor: 0,
			}
		default:
			panic("can't create list model")
	}

}

func (m playerListModel) Init() tea.Cmd {
	return nil
}

func (m playerListModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			homeModel := NewHomeModel()
			return RootScreen().SwitchScreen(homeModel)
		case tea.KeyUp:
			if m.cursor > 0 {
					m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
					m.cursor++
			}
		}
	}
	return m, nil
}

func (m playerListModel) View() string {
	var list string

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
				cursor = ">"
		}
		list += fmt.Sprintf("%s %s\n", cursor, fmt.Sprintf("%s %s", choice.FirstName, choice.LastName))
	}

	player := m.choices[m.cursor]

	displayedPlayer := fmt.Sprintf(
		`%s %s
%s %s`,
		headingStyle.Width(15).Render("First Name"),
		headingStyle.Width(15).Render("Last Name"),
		textStyle.Width(15).Render(player.FirstName),
		textStyle.Width(15).Render(player.LastName),
	)

	return lipgloss.JoinHorizontal(lipgloss.Top, listBox.Render(list), displayedPlayer)
}
