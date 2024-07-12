package tui

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ls13g12/hockey-app/src/pkg/common"
	"go.mongodb.org/mongo-driver/mongo"
)

var globalDB *mongo.Database

func NewTuiApp(cfg common.TuiAppConfig, logger *slog.Logger, db *mongo.Database) {
	if err := db.Client().Ping(context.TODO(), nil); err != nil {
		fmt.Printf("Error pinging db: %v", err)
		os.Exit(1)
	} else {
		globalDB = db
	}

	p := tea.NewProgram(RootScreen(db), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

type rootScreenModel struct {
	model  tea.Model
}

func RootScreen(db *mongo.Database) rootScreenModel {
	var rootModel tea.Model

	screen_one := ScreenOne(db)
	rootModel = &screen_one

	return rootScreenModel{
			model: rootModel,
	}
}

func (m rootScreenModel) Init() tea.Cmd {
	return m.model.Init()
}

func (m rootScreenModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m.model.Update(msg)
}

func (m rootScreenModel) View() string {
	return m.model.View()
}

func (m rootScreenModel) SwitchScreen(model tea.Model) (tea.Model, tea.Cmd) {
	m.model = model
	return m.model, m.model.Init()
}
