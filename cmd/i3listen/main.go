package main

import (
	"encoding/json"
	"log"
	"os"

	"go.i3wm.org/i3"
)

func main() {
	var types []i3.EventType
	for _, arg := range os.Args[1:] {
		t := i3.EventType(arg)
		switch t {
		case i3.WorkspaceEventType, i3.OutputEventType, i3.ModeEventType, i3.WindowEventType, i3.TickEventType:
			types = append(types, t)
		default:
			log.Fatalf("invalid event type %s", arg)
		}
	}

	if len(types) == 0 {
		log.Fatal("no event types given")
	}

	er := i3.Subscribe(types...)
	defer er.Close()

	enc := json.NewEncoder(os.Stdout)
	for er.Next() {
		if err := enc.Encode(er.Event()); err != nil {
			log.Fatal(err)
		}
	}
}
