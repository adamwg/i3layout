package main

import (
	"os"

	"github.com/adamwg/i3layout"
	"github.com/adamwg/i3layout/server"
	"go.i3wm.org/i3"
	"gopkg.in/alecthomas/kingpin.v2"
)

func main() {
	app := kingpin.New("i3layout", "layout utilities for i3")

	app.Command("serve", "run the layout server").Action(serveCommand)
	addOneshotCommand(app)
	addSetModeCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func addOneshotCommand(app *kingpin.Application) {
	cmd := app.Command("oneshot", "lay out a workspace now")
	wsName := cmd.Arg("workspace", "workspace to lay out; defaults to the focused workspace").
		Default("").String()

	cmd.Action(func(*kingpin.ParseContext) error {
		var ws *i3.Node
		var err error

		if *wsName == "" {
			ws, err = i3layout.GetFocusedWorkspace()
		} else {
			ws, err = i3layout.GetWorkspace(*wsName)
		}
		if err != nil {
			return err
		}

		windows := i3layout.GetWindows(ws)
		if err != nil {
			return err
		}
		template := i3layout.MakeTemplate(windows)
		return template.Apply(ws, windows)
	})
}

func serveCommand(*kingpin.ParseContext) error {
	return server.Serve()
}
