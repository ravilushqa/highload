//go:build tools

package tools

import (
	_ "github.com/golangci/golangci-lint/cmd/golangci-lint"
)

//go:generate go build -v -o ../bin/ github.com/golangci/golangci-lint/cmd/golangci-lint
