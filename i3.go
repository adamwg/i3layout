package i3layout

import (
	"errors"

	"go.i3wm.org/i3"
)

func RunI3Cmd(cmd string) error {
	result, err := i3.RunCommand(cmd)
	if err != nil {
		return err
	}
	for _, r := range result {
		if !r.Success {
			return errors.New(r.Error)
		}
	}

	return nil
}
