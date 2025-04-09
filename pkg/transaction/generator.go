//go:build tools
// +build tools

package main

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/deepmap/oapi-codegen/v2/cmd/oapi-codegen -config=oapi.yaml ./oapi-transaction.yaml
