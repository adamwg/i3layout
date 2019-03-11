package i3layout

import (
	"fmt"
	"time"

	"go.i3wm.org/i3"
)

type layoutFunc func(string, []*i3.Node) *Template

const tempWorkspacePrefix = "i3layout-temp-"

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
	return fn(tempWorkspaceName(), windows)
}

func marks(n *i3.Node) []string {
	return []string{fmt.Sprintf("i3layout-%d", n.Window)}
}

func tempWorkspaceName() string {
	return fmt.Sprintf("%s-%d", tempWorkspacePrefix, time.Now().UnixNano())
}
