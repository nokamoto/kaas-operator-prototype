package main

import (
	"fmt"
	"os"

	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

type Build mg.Namespace

// ControllerGenCRD generates the Custom Resource Definitions (CRDs) for the project.
func (Build) ControllerGenCRD() error {
	return sh.RunV("controller-gen", "crd", "paths=./...", "output:crd:dir=config/crd")
}

// ControllerGenObject generates the deep copy files for the project.
func (Build) ControllerGenObject() error {
	return sh.RunV("controller-gen", "object:headerFile=hack/boilerplate.go.txt", "paths=./...")
}

// ControllerGenRBAC generates the RBAC roles for the project.
func (Build) ControllerGenRBAC() error {
	return sh.RunV("controller-gen", "rbac:roleName=manager-role", "paths=./...", "output:rbac:dir=./config/rbac")
}

// Buf generates the go code from the protobuf definitions using Buf.
func (Build) Buf() error {
	return sh.RunV("buf", "generate")
}

// Mock generates the mock for gRPC client interfaces.
func (Build) Mock() error {
	s, err := sh.Output("mockgen", "./pkg/api/proto/v1alpha1/v1alpha1connect", "ClusterServiceClient")
	if err != nil {
		return fmt.Errorf("failed to generate mock: %w", err)
	}
	if err := os.WriteFile("internal/mock/mock_v1alpha1connect/cluster.go", []byte(s), 0o644); err != nil {
		return fmt.Errorf("failed to write mock file: %w", err)
	}
	return nil
}

// InstallMCPServer installs the MCP server binary.
func (Build) InstallMCPServer() error {
	return sh.RunV("go", "install", "./cmd/mcpserver")
}
