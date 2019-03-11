package i3layout

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"go.i3wm.org/i3"
)

type Node struct {
	Name    string      `json:"name,omitempty"`
	Type    i3.NodeType `json:"type"`
	Layout  i3.Layout   `json:"layout,omitempty"`
	Percent float64     `json:"percent,omitempty"`
	Marks   []string    `json:"marks,omitempty"`
	Nodes   []*Node     `json:"nodes,omitempty"`
}

type Template Node

func (t *Template) Marshal(w io.Writer) error {
	return json.NewEncoder(w).Encode(t)
}

func (t *Template) Apply(ws *i3.Node, windows Windows) error {
	if t == nil {
		return nil
	}

	tmpLayout, err := ioutil.TempFile("", "i3-layout-*.json")
	tmpLayoutPath := tmpLayout.Name()
	if err != nil {
		return err
	}
	t.Marshal(tmpLayout)
	tmpLayout.Close()
	defer os.Remove(tmpLayoutPath)

	if err := RunI3Cmd(fmt.Sprintf("append_layout %s", tmpLayoutPath)); err != nil {
		return err
	}

	var cmds []string
	for _, n := range windows {
		cmds = append(cmds,
			fmt.Sprintf("[con_id=%d] swap container with mark i3layout-%d", n.ID, n.Window),
		)
	}
	cmds = append(cmds, fmt.Sprintf(`rename workspace %s to %s`, ws.Name, killWorkspaceName))
	cmds = append(cmds, fmt.Sprintf(`rename workspace %s to %s`, tempWorkspaceName, ws.Name))
	if ws.Focused {
		cmds = append(cmds, fmt.Sprintf(`workspace %s`, ws.Name))
	}
	cmds = append(cmds, `[con_mark="^i3layout-"] kill`)

	cmd := strings.Join(cmds, "; ")
	if err := RunI3Cmd(cmd); err != nil {
		return err
	}

	return nil
}
