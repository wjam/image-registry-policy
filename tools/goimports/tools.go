// +build tools

package tools

import (
	_ "golang.org/x/tools/cmd/goimports"
)

//go:generate go build -v -o=../../bin/goimports golang.org/x/tools/cmd/goimports
