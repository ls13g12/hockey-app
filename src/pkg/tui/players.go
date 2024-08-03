package tui

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/lipgloss"
	"github.com/ls13g12/hockey-app/src/pkg/db"
)

const (
	firstName = iota
	lastName
	homeShirtNumber
)

var (
	playerBox      = lipgloss.NewStyle().Width(80).Padding(0, 2)
)

func generatePlayerCardString(player db.Player) string {
	displayedPlayer := fmt.Sprintf(
		`%s %s
%s %s

%s %s
`,
		headingStyle.Width(15).Render("First Name"),
		headingStyle.Width(15).Render("Last Name"),
		textStyle.Width(15).Render(player.FirstName),
		textStyle.Width(15).Render(player.LastName),
		headingStyle.Width(15).Render("Home Number"),
		headingStyle.Width(15).Render("Away Number"),
	)
	if player.HomeShirtNumber > 0 {
		displayedPlayer += textStyle.Width(15).Render(strconv.Itoa(player.HomeShirtNumber))
	} else {
		displayedPlayer += textStyle.Width(15).Render(" ")
	}

	if player.AwayShirtNumber > 0 {
		displayedPlayer += textStyle.Width(15).Render(strconv.Itoa(player.AwayShirtNumber))
	} else {
		displayedPlayer += textStyle.Width(15).Render(" ")
	}
	return displayedPlayer
}

func generatePlayerInputs() ([]textinput.Model, error) {
	var inputs []textinput.Model = make([]textinput.Model, 2)
	inputs[firstName] = textinput.New()
	inputs[firstName].Focus()
	inputs[firstName].CharLimit = 20
	inputs[firstName].Width = 20
	inputs[firstName].Prompt = ""
	inputs[firstName].Validate = nameValidator

	inputs[lastName] = textinput.New()
	inputs[lastName].CharLimit = 30
	inputs[lastName].Width = 30
	inputs[lastName].Prompt = ""
	inputs[lastName].Validate = nameValidator

	return inputs, nil
}

func nameValidator(s string) error {
	s = strings.ReplaceAll(s, " ", "")

	isAlpha := regexp.MustCompile(`[A-Za-z]+$`).MatchString(s)
	if !isAlpha {
		return fmt.Errorf("names may only contain letters")
	}

	if len(s) > 10 {
		return fmt.Errorf("names must be less than 20 characters")
	}

	if len(s) < 1 {
		return fmt.Errorf("names cannot be empty")
	}

	return nil
}
