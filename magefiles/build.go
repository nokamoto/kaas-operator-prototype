package main

import (
	_ "embed"
	"fmt"
	"os"
	"text/template"

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
	for _, c := range controllers {
		rbac := fmt.Sprintf("rbac:roleName=%s-manager-role,fileName=%s.yaml", c, c)
		path := fmt.Sprintf("paths=./internal/controller/%s/...", c)
		if err := sh.RunV("controller-gen", rbac, path, "output:rbac:dir=./config/rbac"); err != nil {
			return fmt.Errorf("failed to generate RBAC for %s: %w", c, err)
		}
	}
	return nil
}

//go:embed templates/manager.yaml.tmpl
var managerYAMLTemplate string

// ManagerYAML generates the manager YAML file for the project.
func (Build) ManagerYAML() error {
	for _, c := range controllers {
		t, err := template.New("manager.yaml").Parse(managerYAMLTemplate)
		if err != nil {
			return fmt.Errorf("failed to parse manager YAML template: %w", err)
		}
		file, err := os.Create(fmt.Sprintf("config/manager/%s.yaml", c))
		if err != nil {
			return fmt.Errorf("failed to create manager YAML file for %s: %w", c, err)
		}
		defer file.Close()
		if err := t.Execute(file, c); err != nil {
			return fmt.Errorf("failed to execute manager YAML template for %s: %w", c, err)
		}
	}
	return nil
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
