package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// ControllerGenCRD generates the Custom Resource Definitions (CRDs) for the project.
func (Build) ControllerGenCRD() error {
	return sh.RunV("controller-gen", "crd", "paths=api/...", "output:crd:dir=./config/crd")
}
