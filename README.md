# Planner TUI

A terminal UI todo list with a week calendar view, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **Week calendar view** — default view showing 7 days with tasks under each day
- **Day cursor** — `←`/`→` moves between columns, wraps to next/previous week
- **Todo list** — `t` switches to list view with `↑`/`↓` navigation
- **Toggle done** — `Enter` marks tasks complete/incomplete
- **Auto-save** — tasks persist to `~/.todos.json`
- **Visual cues** — today highlighted in green, overdue tasks in red, completed with ✓

## Key Bindings

| Key | View | Action |
|-----|------|--------|
| `←` `→` | Week | Move day cursor |
| `t` | Any | Toggle week / list view |
| `↑` `k` | List | Move cursor up |
| `↓` `j` | List | Move cursor down |
| `Enter` | List | Toggle task done |
| `q` `Ctrl+C` | Any | Quit |

## Install

```bash
git clone https://github.com/Momokh99/planner-tui
cd planner-tui
go build -o planner-tui
./planner-tui
```

Or run without building:

```bash
go run .
```

## Hyprland Integration

### Launch with a keybind

```conf
# ~/.config/hypr/hyprland.conf
bind = SUPER, T, exec, kitty -e /path/to/planner-tui
```

### Scratchpad (toggle with same key)

```conf
bind = SUPER, T, exec, kitty --class planner-tui -e /path/to/planner-tui
windowrule = float, ^(planner-tui)$
windowrule = workspace special silent, ^(planner-tui)$
bind = SUPER SHIFT, T, togglespecialworkspace
```

### Auto-start on login

```conf
exec-once = kitty --class planner-tui -e /path/to/planner-tui
```

## Data

Tasks are stored in `~/.todos.json`. You can edit it manually:

```json
[
  {
    "id": 1,
    "title": "Buy groceries",
    "completed": false,
    "due_date": "2026-06-19T00:00:00Z"
  }
]
```

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling
