//go:build tools

package tools

import (
	_ "github.com/bufbuild/buf/cmd/buf"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
