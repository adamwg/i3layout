package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync/atomic"

	"github.com/adamwg/i3layout"
	"go.i3wm.org/i3"
)

var (
	workspaceLayouts = map[string]string{}
)

func Serve() (err error) {
	// We set processing to true when handling an event so that we can skip
	// handling the events that the event handling causes. This isn't ideal, but
	// there's no reliable way to filter the events and avoid getting stuck in
	// an infinite loop of our own creation.
	var processing atomic.Value
	processing.Store(false)

	er := i3.Subscribe(
		i3.WindowEventType,
		i3.TickEventType,
	)
	defer func() {
		if err2 := er.Close(); err2 != nil {
			err = err2
		}
	}()
	for er.Next() && err == nil {
		if processing.Load().(bool) {
			continue
		}
		processing.Store(true)

		event := er.Event()
		switch ev := event.(type) {
		case *i3.WindowEvent:
			// This is done in a goroutine so that we can skip over the events
			// it causes.
			go func() {
				err = handleWindowEvent(ev)
				processing.Store(false)
			}()
		case *i3.TickEvent:
			// This is done in a goroutine so that we can skip over the events
			// it causes.
			go func() {
				err = handleTickEvent(ev)
				processing.Store(false)
			}()
		}
	}

	return err
}

func handleWindowEvent(ev *i3.WindowEvent) error {
	if ev.Container.Window == 0 {
		return nil
	}

	var ws *i3.Node
	var err error

	switch ev.Change {
	case "new", "move":
		log.Printf("handling %q event for window %q", ev.Change, ev.Container.Name)
		ws, err = i3layout.GetWorkspaceForWindow(&ev.Container)
		if err != nil {
			return err
		}

		if i3layout.IsFocusedWorkspace(ws) {
			ws.Focused = true
		}

	case "close":
		log.Printf("handling %q event for window %q", ev.Change, ev.Container.Name)
		ws, err = i3layout.GetFocusedWorkspace()
		if err != nil {
			return err
		}

	default:
		return nil
	}

	if workspaceLayouts[ws.Name] == "" {
		workspaceLayouts[ws.Name] = i3layout.LayoutNames()[0]
	}

	windows := i3layout.GetWindows(ws)
	template := i3layout.MakeTemplate(workspaceLayouts[ws.Name], windows)

	if err := template.Apply(ws, windows); err != nil {
		return err
	}

	if ev.Change == "new" {
		return i3layout.RunI3Cmd(fmt.Sprintf("[con_id=%d] focus", ev.Container.ID))
	}

	return nil
}

type ChangeLayoutMessage struct {
	Tag           string `json:"tag"`
	WorkspaceName string `json:"workspace_name"`
	LayoutName    string `json:"layout_name"`
}

const ChangeLayoutTag = "i3layout-change-layout"

func handleTickEvent(ev *i3.TickEvent) error {
	var payload ChangeLayoutMessage
	err := json.Unmarshal([]byte(ev.Payload), &payload)
	if err != nil || payload.Tag != ChangeLayoutTag {
		// The tick payload wasn't from i3layout, so ignore it.
		return nil
	}

	ws, err := i3layout.GetWorkspace(payload.WorkspaceName)
	if err != nil {
		return err
	}

	layoutName := payload.LayoutName
	if payload.LayoutName == "next" {
		layoutName = findOffsetLayout(payload.WorkspaceName, 1)
	} else if payload.LayoutName == "prev" {
		layoutName = findOffsetLayout(payload.WorkspaceName, -1)
	} else if !i3layout.LayoutExists(payload.LayoutName) {
		return errors.New("invalid layout name")
	}

	log.Printf("changing layout for workspace %q to %q", ws.Name, layoutName)

	workspaceLayouts[ws.Name] = layoutName
	windows := i3layout.GetWindows(ws)
	template := i3layout.MakeTemplate(layoutName, windows)

	if i3layout.IsFocusedWorkspace(ws) {
		ws.Focused = true
	}

	if err := template.Apply(ws, windows); err != nil {
		return err
	}

	return nil
}

func findOffsetLayout(wsName string, offset int) string {
	layoutNames := i3layout.LayoutNames()
	currentLayout, ok := workspaceLayouts[wsName]
	if !ok {
		return layoutNames[0]
	}

	// Make offset positive for use in arithmetic below.
	for offset < 0 {
		offset += len(layoutNames)
	}
	for i, n := range layoutNames {
		if n == currentLayout {
			return layoutNames[(i+offset)%len(layoutNames)]
		}
	}

	return layoutNames[0]
}
