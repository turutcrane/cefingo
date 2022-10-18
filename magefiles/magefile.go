//go:build mage
// +build mage

package main

import (
	"github.com/magefile/mage/sh"
)

// Build Capi
func Capi() error {
	return sh.RunV("go", "run", "github.com/turutcrane/gencefingo@latest", "-pkgdir", ".")
}
