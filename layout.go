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

const (
	tallWidth = float64(296875.0 / 1000000.0)
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
	cmds = append(cmds, fmt.Sprintf(`rename workspace %s to i3layout-kill`, ws.Name))
	cmds = append(cmds, fmt.Sprintf(`rename workspace i3layout-temp to %s`, ws.Name))
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

func MakeTemplate(windows []*i3.Node) *Template {
	if len(windows) == 0 {
		return nil
	}

	layout := &Template{
		Name:  "i3layout-temp",
		Type:  i3.WorkspaceNode,
		Marks: []string{"i3layout-temp"},
		Nodes: []*Node{
			{
				Type:  i3.Con,
				Marks: marks(windows[0]),
			},
		},
	}

	if len(windows) == 1 {
		return layout
	}

	layout.Nodes[0].Percent = tallWidth
	layout.Nodes = append(layout.Nodes, &Node{
		Type:    i3.Con,
		Percent: 1.0 - tallWidth,
	})

	if len(windows) == 2 {
		layout.Nodes[1].Layout = i3.Tabbed
		layout.Nodes[1].Marks = marks(windows[1])
		return layout
	}

	layout.Nodes[1].Layout = i3.SplitV
	for _, n := range windows[1:] {
		layout.Nodes[1].Nodes = append(layout.Nodes[1].Nodes, &Node{
			Type:   i3.Con,
			Layout: i3.Tabbed,
			Marks:  marks(n),
		})
	}

	return layout
}

func marks(n *i3.Node) []string {
	return []string{fmt.Sprintf("i3layout-%d", n.Window)}
}
