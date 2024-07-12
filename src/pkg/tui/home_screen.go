package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ls13g12/hockey-app/src/pkg/db"
	"go.mongodb.org/mongo-driver/mongo"
)

type screenOneModel struct {
	db *mongo.Database
	choices []string
	cursor int
}

func ScreenOne(db *mongo.Database) screenOneModel {
	return screenOneModel{
		db: db,
		choices: []string{"players", "matches"},
	}
}

func (m screenOneModel) Init() tea.Cmd {
	return nil
}


func (m screenOneModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
			switch msg.Type {
					case tea.KeyCtrlC:
							return m, tea.Quit

					case tea.KeyUp:
							if m.cursor > 0 {
									m.cursor--
							}

					case tea.KeyDown:
							if m.cursor < len(m.choices)-1 {
									m.cursor++
							}

					case tea.KeyEnter, tea.KeyBackspace:
							switch m.cursor {
								case 0:
									players, _ := db.AllPlayers(m.db)
									sidebarModel := playerSidebarModel{
										db: m.db,
										choices: players,
										cursor: 0,
									}
									createModel := playerCreateModel{
										inputs: generateInputs(),
										focused: 0,
										err: nil,
									}
									screen_two := MainScreen(sidebarModel, createModel)
									return RootScreen(m.db).SwitchScreen(&screen_two)
								// case 1:
								// 	screen_two := MainScreen()
								// 	return RootScreen().SwitchScreen(&screen_two)
							}
			}
	}
	return m, nil
}

func (m screenOneModel) View() string {
	var s string
	for i, choice := range m.choices {

		cursor := " "
		if m.cursor == i {
				cursor = ">" 
		}

		// Render the row
		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s
}
