package main

import (
	"fmt"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

var Default = All

// Install installs the necessary tools for the project.
func Install() error {
	tools := []string{
		"sigs.k8s.io/controller-tools/cmd/controller-gen",
		"golang.org/x/tools/cmd/goimports",
		"mvdan.cc/gofumpt",
		"github.com/google/ko",
	}
	for _, tool := range tools {
		if err := sh.Run("go", "install", fmt.Sprintf("%s@latest", tool)); err != nil {
			return err
		}
	}
	return nil
}

// All runs all the necessary tasks to prepare the project.
func All() {
	mg.SerialDeps(
		Install,
		Build.ControllerGenCRD,
		Build.ControllerGenObject,
		Format.Go,
	)
}
