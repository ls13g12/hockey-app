package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ls13g12/hockey-app/src/pkg/db"
)

type DeleteModel interface {
	Init() tea.Cmd
	View() string
	Update(msg tea.Msg) (DeleteModel, tea.Cmd)
}

type playerDeleteModel struct {
	confirmChoice bool
	choices []db.Player
	cursor  int
}

func NewDeleteModel(resource string) (tea.Model, error) {
	switch resource {
	case "players":
		players, err := db.AllPlayers(globalDB)
		if err != nil {
			return nil, fmt.Errorf("error fetching players")
		}
		return playerDeleteModel{
			choices: players,
			cursor:  0,
		}, nil
	default:
		return nil, fmt.Errorf("model for list action not implemented on resource %s", resource)
	}
}

func (m playerDeleteModel) Init() tea.Cmd {
	return nil
}

func (m playerDeleteModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.confirmChoice {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyCtrlC, tea.KeyEsc:
				m.confirmChoice = false
			case tea.KeyEnter:
				player := m.choices[m.cursor]
				err := db.DeletePlayer(globalDB, player.PlayerID)
				if err != nil {
					panic(err)
				}
				homeModel := NewHomeModel("player successfully deleted")
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
		case tea.KeyUp:
			if m.cursor > 0 {
				m.cursor--
			}
		case tea.KeyDown:
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case tea.KeyEnter:
			m.confirmChoice = true
		}
	}
	return m, nil
}

func (m playerDeleteModel) View() string {
	var list string

	if len(m.choices) == 0 {
		return notificationBox.Render("No players to show")
	}

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = ">"
		}
		list += fmt.Sprintf("%s %s\n", cursor, fmt.Sprintf("%s %s", choice.FirstName, choice.LastName))
	}

	player := m.choices[m.cursor]

	displayedPlayer := generatePlayerCardString(player)

	var confirmMessage string
	if m.confirmChoice {
		confirmMessage = "Click Enter again to confirm"
		return lipgloss.JoinVertical(lipgloss.Top, notificationBox.Render(confirmMessage), lipgloss.JoinHorizontal(lipgloss.Top, listBox.Render(list), displayedPlayer))
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, listBox.Render(list), playerBox.Render(displayedPlayer))
}
