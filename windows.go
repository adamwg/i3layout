package i3layout

import (
	"errors"

	"go.i3wm.org/i3"
)

type Windows []*i3.Node

func GetWindows(n *i3.Node) Windows {
	var windows []*i3.Node
	for _, n := range n.Nodes {
		if len(n.Nodes) > 0 {
			windows = append(windows, GetWindows(n)...)
		} else if n.Window != 0 {
			windows = append(windows, n)
		}
	}

	return windows
}

func GetFocusedWorkspace() (*i3.Node, error) {
	tree, err := i3.GetTree()
	if err != nil {
		return nil, err
	}

	ws := tree.Root.FindFocused(func(n *i3.Node) bool { return n.Type == i3.WorkspaceNode })
	ws.Focused = true

	return ws, nil
}

func GetWorkspaceForWindow(window *i3.Node) (*i3.Node, error) {
	tree, err := i3.GetTree()
	if err != nil {
		return nil, err
	}

	return tree.Root.FindChild(workspaceContaining(window)), nil
}

func GetWorkspace(name string) (*i3.Node, error) {
	tree, err := i3.GetTree()
	if err != nil {
		return nil, err
	}

	ws := tree.Root.FindChild(func(n *i3.Node) bool { return n.Type == i3.WorkspaceNode && n.Name == name })

	if ws == nil {
		return nil, errors.New("workspace not found")
	}
	return ws, nil
}

func IsFocusedWorkspace(workspace *i3.Node) bool {
	if workspace.Focused {
		return true
	}

	for _, n := range workspace.Nodes {
		if n.Focused {
			return true
		}
		if IsFocusedWorkspace(n) {
			return true
		}
	}

	return false
}

func workspaceContaining(window *i3.Node) func(*i3.Node) bool {
	return func(n *i3.Node) bool {
		if n.Type != i3.WorkspaceNode {
			return false
		}
		return hasChild(n, window)
	}
}

func hasChild(n, window *i3.Node) bool {
	for _, child := range n.Nodes {
		if child.ID == window.ID {
			return true
		}
		if hasChild(child, window) {
			return true
		}
	}

	return false
}
