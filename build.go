package main

import (
	"bytes"
	"errors"
	"os"
	ex "os/exec"
	"path/filepath"
)

const execName = "main"

func exec(name string, args ...string) (string, error) {
	cmd := ex.Command(name, args...)

	errs := bytes.NewBufferString("")
	out := bytes.NewBufferString("")
	cmd.Stderr = errs
	cmd.Stdout = out
	err := cmd.Run()
	if err != nil {
		return "", errors.New(errs.String())
	}
	return out.String(), nil
}

func build(module string) string {
	execPath := filepath.Join(outdir, execName)

	_, err := exec("go", "build", "-o", execPath, module)
	handle(err)

	sized, err := exec("go", "tool", "nm", "-size", execPath)
	handle(err)

	err = os.Remove(execPath)
	handle(err)

	return sized
}
