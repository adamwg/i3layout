package i3layout

import (
	"fmt"

	"go.i3wm.org/i3"
)

type layoutFunc func([]*i3.Node) *Template

const (
	tempWorkspaceName = "i3layout-temp"
	killWorkspaceName = "i3layout-kill"
)

var (
	layoutFuncs = map[string]layoutFunc{
		"tall":    tallLayout,
		"columns": columnsLayout,
		"rows":    rowsLayout,
		"full":    fullLayout,
	}
)

func LayoutNames() []string {
	return []string{"tall", "columns", "rows", "full"}
}

func LayoutExists(layoutName string) bool {
	_, ok := layoutFuncs[layoutName]
	return ok
}

func MakeTemplate(layoutName string, windows []*i3.Node) *Template {
	if len(windows) == 0 {
		return nil
	}

	fn := layoutFuncs[layoutName]
	return fn(windows)
}

func marks(n *i3.Node) []string {
	return []string{fmt.Sprintf("i3layout-%d", n.Window)}
}
