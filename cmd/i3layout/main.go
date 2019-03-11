package main

import (
	"encoding/json"
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
	addClientCommand(app)

	kingpin.MustParse(app.Parse(os.Args[1:]))
}

func addOneshotCommand(app *kingpin.Application) {
	cmd := app.Command("oneshot", "lay out a workspace now")
	wsName := cmd.Arg("workspace", "workspace to lay out; defaults to the focused workspace").
		Default("").String()
	layoutName := cmd.Flag("layout", "name of the layout to apply").Default("tall").Enum(i3layout.LayoutNames()...)

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
		template := i3layout.MakeTemplate(*layoutName, windows)
		return template.Apply(ws, windows)
	})
}

func addClientCommand(app *kingpin.Application) {
	cmd := app.Command("client", "client commands for the i3layout server")
	changeLayoutCmd := cmd.Command("change-layout", "change the layout for a workspace")
	wsName := changeLayoutCmd.Flag("workspace", "workspace to set layout for; defaults to the focused workspace").
		Default("").String()
	layoutName := changeLayoutCmd.Arg("layout", "name of the layout to apply, or next or prev to cycle").
		Default("next").Enum(append(i3layout.LayoutNames(), "next", "prev")...)

	changeLayoutCmd.Action(func(*kingpin.ParseContext) error {
		if *wsName == "" {
			ws, err := i3layout.GetFocusedWorkspace()
			if err != nil {
				return err
			}
			*wsName = ws.Name
		}

		msg := server.ChangeLayoutMessage{
			Tag:           server.ChangeLayoutTag,
			WorkspaceName: *wsName,
			LayoutName:    *layoutName,
		}
		bs, err := json.Marshal(msg)
		if err != nil {
			return err
		}

		_, err = i3.SendTick(string(bs))
		return err
	})
}

func serveCommand(*kingpin.ParseContext) error {
	return server.Serve()
}
