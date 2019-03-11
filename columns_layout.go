package i3layout

import "go.i3wm.org/i3"

func columnsLayout(windows []*i3.Node) *Template {
	layout := &Template{
		Name:   tempWorkspaceName,
		Type:   i3.WorkspaceNode,
		Layout: i3.SplitH,
	}

	for _, w := range windows {
		layout.Nodes = append(layout.Nodes, &Node{
			Type:  i3.Con,
			Marks: marks(w),
		})
	}

	return layout
}