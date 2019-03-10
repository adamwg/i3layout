package main

import (
	"encoding/json"
	"log"
	"os"

	"github.com/adamwg/i3layout"
	"go.i3wm.org/i3"
)

func main() {
	mode := "focused"
	if len(os.Args) >= 2 {
		mode = os.Args[1]
	}

	var ws *i3.Node
	var err error
	if mode == "focused" {
		ws, err = i3layout.GetFocusedWorkspace()
		if err != nil {
			log.Fatal(err)
		}
	} else {
		tree, err := i3.GetTree()
		if err != nil {
			log.Fatal(err)
		}
		ws = tree.Root
	}

	if mode == "focused" || mode == "all" {
		windows := i3layout.GetWindows(ws)
		if err := json.NewEncoder(os.Stdout).Encode(windows); err != nil {
			log.Fatal(err)
		}
	} else {
		if err := json.NewEncoder(os.Stdout).Encode(ws); err != nil {
			log.Fatal(err)
		}
	}
}
