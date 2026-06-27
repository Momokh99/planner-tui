# Planner TUI

A terminal UI todo list with a week calendar view, built with [Bubble Tea](https://github.com/charmbracelet/bubbletea).

## Features

- **Week calendar view** — default view showing 7 days with tasks under each day
- **Day cursor** — `←`/`→` moves between columns, wraps to next/previous week
- **Todo list** — `t` switches to list view with `↑`/`↓` navigation
- **Add / Edit / Delete** — `a` add, `e` edit, `d` delete with confirmation
- **Toggle done** — `Enter` marks tasks complete/incomplete
- **Auto-save** — tasks persist to `~/.todo.json`
- **Visual cues** — today highlighted in green, cursor column in gray, overdue tasks in red, completed with ✓
- **Task details** — full untruncated todos for cursor day shown below the week grid
- **Robust IDs** — auto-incrementing counter, no ID reuse

## Key Bindings

| Key | View | Action |
|-----|------|--------|
| `←` `→` | Week | Move day cursor |
| `g` | Week | Jump cursor to today |
| `t` | Any | Toggle week / list view |
| `↑` `k` | List | Move cursor up |
| `↓` `j` | List | Move cursor down |
| `Enter` | List | Toggle task done |
| `a` | List | Add new task |
| `e` | List | Edit selected task |
| `d` | List | Delete selected task |
| `y` | Confirm | Confirm delete |
| `n` `Esc` | Confirm | Cancel delete |
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

Tasks are stored in `~/.todo.json`. You can edit it manually:

```json
[
  {
    "ID": 1,
    "Title": "Buy groceries",
    "Completed": false,
    "DueDate": "2026-06-19T00:00:00Z"
  }
]
```

## Built With

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — styling
