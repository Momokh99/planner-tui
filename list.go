package main

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("57"))
var overdueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

func (m model) listView() string {
	s := "My Todos\n\n"
	if len(m.todos) == 0 {
		s += "No todos yet.\n"
	}
	for i, t := range m.todos {
		cursor := "  "
		if m.cursor == i {
			cursor = "▸ "
		}
		check := "○"
		if t.Completed {
			check = "●"
		}
		date := t.DueDate.Format("Jan 02")

		line := cursor + check + " " + t.Title + " — " + date
		if !t.Completed && t.DueDate.Before(time.Now()) {
			line = overdueStyle.Render(line)
		}
		if m.cursor == i {
			line = selectedStyle.Render(line)
		}
		s += line + "\n"
	}
	s += "\n  [a] add  [e] edit  [d] delete  [Enter] toggle  [t] week  [q] quit"
	return s
}
