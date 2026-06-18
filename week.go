package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var todayHeaderStyle = lipgloss.NewStyle().
	Width(14).Align(lipgloss.Center).Bold(true).
	Background(lipgloss.Color("58")) // dark green
var todayBoxStyle = lipgloss.NewStyle().
	Width(14).Align(lipgloss.Center).
	Background(lipgloss.Color("58"))
var dayHeaderStyle = lipgloss.NewStyle().Width(14).Align(lipgloss.Center).Bold(true)
var dayBoxStyle = lipgloss.NewStyle().Width(14).Align(lipgloss.Center)
var overdueStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196")) // red

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) > max {
		return string(runes[:max-1]) + "…"
	}
	return s
}

func todosForDay(todos []Todo, date time.Time) []Todo {
	var result []Todo
	for _, t := range todos {
		y1, m1, d1 := t.DueDate.Date()
		y2, m2, d2 := date.Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			result = append(result, t)
		}
	}
	return result
}

func mondayOfWeek(t time.Time) time.Time {
	offset := int(t.Weekday()) - 1
	if t.Weekday() == time.Sunday {
		offset = -6
	}
	return t.AddDate(0, 0, -offset)
}

func (m model) weekView() string {
	// Build header row (Mon 16, Tue 17, ...)
	var headers []string
	var columns []string
	for i := 0; i < 7; i++ {
		d := m.weekStart.AddDate(0, 0, i)
		isToday := d.Year() == time.Now().Year() && d.YearDay() == time.Now().YearDay()
		headerStyle := dayHeaderStyle
		boxStyle := dayBoxStyle
		if isToday {
			headerStyle = todayHeaderStyle
			boxStyle = todayBoxStyle
		}
		header := d.Format("Mon 2")
		headers = append(headers, headerStyle.Render(header))
		// Build content for that day
		content := ""
		for _, t := range todosForDay(m.todos, d) {
			title := truncate(t.Title, 11)
			if t.Completed {
				content += "✓ " + title + "\n"
			} else if t.DueDate.Before(time.Now()) {
				content += overdueStyle.Render("⚠ "+title) + "\n"
			} else {
				content += "  " + title + "\n"
			}
		}
		columns = append(columns, boxStyle.Render(content))
	}

	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headers...)
	contentRow := lipgloss.JoinHorizontal(lipgloss.Top, columns...)

	// Month/year title + week navigation hint
	title := m.weekStart.Format("January 2006")
	var markers []string
	for i := 0; i < 7; i++ {
		if i == m.dayCursor {
			markers = append(markers, centerText("▼", 14))
		} else {
			markers = append(markers, centerText("", 14))
		}
	}
	markerRow := lipgloss.JoinHorizontal(lipgloss.Top, markers...)
	return title + "\n" + markerRow + "\n" + headerRow + "\n" + contentRow + "\n\n  ←  → day  [t] list  [q] quit"

}
func centerText(s string, width int) string {
	if len(s) >= width {
		return s
	}
	left := (width - len([]rune(s))) / 2
	right := width - left - len([]rune(s))
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}
