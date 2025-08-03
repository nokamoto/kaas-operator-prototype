package main

import (
	"github.com/magefile/mage/mg"
	"github.com/magefile/mage/sh"
)

// Python provides Mage targets for Python environment setup and package installation.
type Python mg.Namespace

// Venv creates a Python virtual environment in the .venv directory.
func (Python) Venv() error {
	return sh.Run("python3", "-m", "venv", ".venv")
}

// Install installs the 'uv' package using pip.
func (Python) Install() error {
	return sh.Run("pip", "install", "uv")
}
