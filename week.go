package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

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
	cw := columnWidth(m.width, 7)

	dayHeader := lipgloss.NewStyle().Width(cw).Align(lipgloss.Center).Bold(true)
	todayHeader := dayHeader.Background(lipgloss.Color("58"))
	dayBox := lipgloss.NewStyle().Width(cw).Align(lipgloss.Center)
	todayBox := dayBox.Background(lipgloss.Color("58"))
	cursorBox := dayBox.Background(lipgloss.Color("239"))
	overdue := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	maxLines := (m.height - 10) / 2
	if maxLines < 2 {
		maxLines = 2
	}

	var headers []string
	var columns []string
	for i := 0; i < 7; i++ {
		d := m.weekStart.AddDate(0, 0, i)
		isToday := d.Year() == time.Now().Year() && d.YearDay() == time.Now().YearDay()
		headerStyle := dayHeader
		boxStyle := dayBox
		if isToday {
			headerStyle = todayHeader
			boxStyle = todayBox
		} else if i == m.dayCursor {
			headerStyle = dayHeader.Background(lipgloss.Color("239"))
			boxStyle = cursorBox
		}
		header := d.Format("Mon 2")
		headers = append(headers, headerStyle.Render(header))
		content := ""
		dayTasks := todosForDay(m.todos, d)
		for ti, t := range dayTasks {
			if ti >= maxLines {
				remaining := len(dayTasks) - ti
				content += fmt.Sprintf("  +%d more\n", remaining)
				break
			}
			title := truncate(t.Title, cw-3)
			if t.Completed {
				content += "✓ " + title + "\n"
			} else if t.DueDate.Before(time.Now()) {
				content += overdue.Render("⚠ "+title) + "\n"
			} else {
				content += "  " + title + "\n"
			}
		}
		columns = append(columns, boxStyle.Render(content))
	}

	headerRow := lipgloss.JoinHorizontal(lipgloss.Top, headers...)
	contentRow := lipgloss.JoinHorizontal(lipgloss.Top, columns...)
	title := m.weekStart.Format("January 2006")
	var markers []string
	for i := 0; i < 7; i++ {
		if i == m.dayCursor {
			markers = append(markers, centerText("▼", cw))
		} else {
			markers = append(markers, centerText("", cw))
		}
	}
	markerRow := lipgloss.JoinHorizontal(lipgloss.Top, markers...)

	divider := strings.Repeat("─", m.width) + "\n"
	cursorDate := m.weekStart.AddDate(0, 0, m.dayCursor)
	dayTodos := todosForDay(m.todos, cursorDate)

	details := "\n" + divider
	details += "  " + cursorDate.Format("Mon Jan 2, 2006") + "\n\n"
	if len(dayTodos) == 0 {
		details += "  No tasks\n"
	}
	for _, t := range dayTodos {
		check := "○"
		if t.Completed {
			check = "●"
		}
		line := "  " + check + " " + t.Title + " — " + t.DueDate.Format("Jan 02")
		if !t.Completed && t.DueDate.Before(time.Now()) {
			line = overdue.Render(line)
		}
		details += line + "\n"
	}

	return title + "\n" + markerRow + "\n" + headerRow + "\n" + contentRow + details + "\n  ←  → day  [g] today  [t] list  [q] quit"
}
func centerText(s string, width int) string {
	if len(s) >= width {
		return s
	}
	left := (width - len([]rune(s))) / 2
	right := width - left - len([]rune(s))
	return strings.Repeat(" ", left) + s + strings.Repeat(" ", right)
}

func columnWidth(totalWidth, columns int) int {
	w := (totalWidth - columns - 1) / columns
	if w < 10 {
		w = 10
	}
	return w
}
