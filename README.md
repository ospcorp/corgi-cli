```
 ██████╗ ██████╗ ██████╗  ██████╗  ██╗     ██████╗██╗     ██╗
██╔════╝██╔═══██╗██╔══██╗██╔════╝  ██║    ██╔════╝██║     ██║
██║     ██║   ██║██████╔╝██║  ███╗ ██║    ██║     ██║     ██║
██║     ██║   ██║██╔══██╗██║   ██║ ██║    ██║     ██║     ██║
╚██████╗╚██████╔╝██║  ██║╚██████╔╝ ██║    ╚██████╗███████╗██║
 ╚═════╝ ╚═════╝ ╚═╝  ╚═╝ ╚═════╝  ╚═╝     ╚═════╝╚══════╝╚═╝
```
A terminal IRC client written in Go, built with the [Charm](https://charm.sh/) TUI stack.

> **Status:** Early development — initial TUI scaffolding only. IRC connectivity not yet implemented.

## Features

- [x] Full-screen terminal UI with scrollable message history
- [x] Text input for composing messages
- [ ] Connect to IRC servers and channels
- [ ] Send and receive IRC messages

## Tech stack

- **[Bubbletea](https://charm.land/bubbletea)** — Elm-inspired TUI framework
- **[Bubbles](https://charm.land/bubbles)** — pre-built TUI components (text input, viewport)
- **[Lipgloss](https://charm.land/lipgloss)** — terminal styling

## Running

```sh
go run .
```

Press `Ctrl+C` or `Esc` to quit.
