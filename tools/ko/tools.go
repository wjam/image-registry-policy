// +build tools

package tools

import (
	_ "github.com/google/ko"
)

//go:generate go build -v -o=../../bin/ko github.com/google/ko
