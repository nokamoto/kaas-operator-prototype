package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// ControllerGenCRD generates the Custom Resource Definitions (CRDs) for the project.
func (Build) ControllerGenCRD() error {
	return sh.RunV("controller-gen", "rbac:roleName=manager-role", "crd", "paths=./...", "output:crd:dir=config/crd")
}

// ControllerGenObject generates the deep copy files for the project.
func (Build) ControllerGenObject() error {
	return sh.RunV("controller-gen", "object:headerFile=hack/boilerplate.go.txt", "paths=./...")
}

// ControllerGenRBAC generates the RBAC roles for the project.
func (Build) ControllerGenRBAC() error {
	return sh.RunV("controller-gen", "rbac:roleName=manager-role", "paths=./...", "output:rbac:dir=./config/rbac")
}
