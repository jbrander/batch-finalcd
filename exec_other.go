//go:build !windows
// +build !windows

package main

import (
	"bytes"
	"fmt"
)

func runCommand(exePath string, args []string, stdout *bytes.Buffer) error {
	return fmt.Errorf("not implemented")
}
