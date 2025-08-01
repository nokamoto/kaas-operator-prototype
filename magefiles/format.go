package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Format mg.Namespace

// Go formats the Go code in the project.
func (Format) Go() error {
	cmds := [][]string{
		{"go", "mod", "tidy"},
		{"gofumpt", "-l", "-w", "."},
		{"goimports", "-l", "-w", "."},
	}
	for _, cmd := range cmds {
		if err := sh.RunV(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}
	return nil
}

// Proto formats the Protocol Buffers files in the project.
func (Format) Proto() error {
	return sh.RunV("buf", "format", "-w", ".")
}

// YAML formats the YAML files in the project.
func (Format) YAML() error {
	return sh.RunV("yamlfmt", ".")
}
