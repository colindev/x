package main

import (
	"bytes"
	"os"
	"os/exec"
)

func clientAdapt(args []string) func([]byte) error {

	end := []byte("\n")

	if len(args) == 0 {
		return func(b []byte) error {
			os.Stdout.Write(append(b, end...))
			return nil
		}
	}

	cmdName, cmdArgs := args[0], args[1:]

	return func(b []byte) error {
		cmd := exec.Command(cmdName, cmdArgs...)
		cmd.Stdin = bytes.NewBuffer(append(b, end...))
		cmd.Stdout = os.Stdout
		return cmd.Run()
	}
}
