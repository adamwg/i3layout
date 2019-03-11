package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/adamwg/i3layout"
	"go.i3wm.org/i3"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("i3window", "window utilities for i3")

	addListCommand(app)
	app.Command("tree", "show the i3 window tree").Action(treeCommand)

	addFocusCommand(app)
	addSwapCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func addListCommand(app *kingpin.Application) {
	cmd := app.Command("list", "list windows")
	wsName := cmd.Flag("workspace", "workspace to list windows from; defaults to the focused workspace").
		Default("").String()
	all := cmd.Flag("all", "list windows from all workspaces; takes precedence over the workspace arg").
		Bool()

	cmd.Action(func(*kingpin.ParseContext) error {
		var ws *i3.Node
		var err error
		if *all {
			tree, err := i3.GetTree()
			if err != nil {
				return err
			}
			ws = tree.Root
		} else if *wsName == "" {
			ws, err = i3layout.GetFocusedWorkspace()
		} else {
			ws, err = i3layout.GetWorkspace(*wsName)
		}
		if err != nil {
			return err
		}

		windows := i3layout.GetWindows(ws)
		return json.NewEncoder(os.Stdout).Encode(windows)
	})
}

func treeCommand(*kingpin.ParseContext) error {
	tree, err := i3.GetTree()
	if err != nil {
		return err
	}
	return json.NewEncoder(os.Stdout).Encode(tree.Root)
}

func addFocusCommand(app *kingpin.Application) {
	cmd := app.Command("focus", "focus windows")
	nextCmd := cmd.Command("next", "focus the next window")
	nextCmd.Action(func(*kingpin.ParseContext) error {
		return focusWindowOffset(1)
	})
	prevCmd := cmd.Command("prev", "focus the previous window")
	prevCmd.Action(func(*kingpin.ParseContext) error {
		return focusWindowOffset(-1)
	})
}

func focusWindowOffset(offset int) error {
	ws, err := i3layout.GetFocusedWorkspace()
	if err != nil {
		return err
	}
	windows := i3layout.GetWindows(ws)

	for offset < 0 {
		offset += len(windows)
	}

	for i, w := range windows {
		if w.Focused {
			win := windows[(i+offset)%len(windows)]
			return i3layout.RunI3Cmd(fmt.Sprintf("[con_id=%d] focus", win.ID))
		}
	}

	return errors.New("could not find focused window")
}

func addSwapCommand(app *kingpin.Application) {
	cmd := app.Command("swap", "swap windows")
	nextCmd := cmd.Command("next", "swap with the next window")
	nextCmd.Action(func(*kingpin.ParseContext) error {
		return swapWindowOffset(1)
	})
	prevCmd := cmd.Command("prev", "swap with the previous window")
	prevCmd.Action(func(*kingpin.ParseContext) error {
		return swapWindowOffset(-1)
	})
}

func swapWindowOffset(offset int) error {
	ws, err := i3layout.GetFocusedWorkspace()
	if err != nil {
		return err
	}
	windows := i3layout.GetWindows(ws)

	for offset < 0 {
		offset += len(windows)
	}

	for i, w := range windows {
		if w.Focused {
			win := windows[(i+offset)%len(windows)]
			return i3layout.RunI3Cmd(fmt.Sprintf("swap with con_id %d", win.ID))
		}
	}

	return errors.New("could not find focused window")
}
