//go:build windows
// +build windows

package main

import (
	"bytes"
	"os/exec"
	"syscall"
)

func runCommand(exePath string, args []string, stdout *bytes.Buffer) error {

	cmd := exec.Command(exePath, args...)
	cmd.Stdout = stdout
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow:    true,
		CreationFlags: 0x08000000,
	}
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil

}
