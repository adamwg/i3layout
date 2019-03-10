package server

import (
	"fmt"
	"log"
	"sync/atomic"

	"github.com/adamwg/i3layout"
	"go.i3wm.org/i3"
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

		// This is done in a goroutine so that we can skip over the events it
		// causes.
		go func() {
			ev := er.Event().(*i3.WindowEvent)
			err = handleWindowEvent(ev)
			processing.Store(false)
		}()
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
	case "new", "move", "floating":
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

	windows := i3layout.GetWindows(ws)
	template := i3layout.MakeTemplate(windows)

	if err := template.Apply(ws, windows); err != nil {
		return err
	}

	if ev.Change == "new" {
		return i3layout.RunI3Cmd(fmt.Sprintf("[con_id=%d] focus", ev.Container.ID))
	}

	return nil
}
