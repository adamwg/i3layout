package i3layout

import "go.i3wm.org/i3"

const (
	tallWidth = float64(0.3)
)

func tallLayout(wsName string, windows []*i3.Node) *Template {
	layout := &Template{
		Name:   wsName,
		Type:   i3.WorkspaceNode,
		Layout: i3.SplitH,
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
