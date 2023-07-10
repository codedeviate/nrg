package helper

import (
	"os/exec"
	"path/filepath"
)

func Preview(command *Command) error {
	var cmd *exec.Cmd
	filename := command.Args[0]
	if filename[0] != '/' {
		stack := GetStack()
		filename = filepath.Join(stack.ActivePath, filename)
	}
	//	cmd = exec.Command("cmd", "/C", "start", command.Args[0])
	cmd = exec.Command("open", "-a", "Preview.app", filename)
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
