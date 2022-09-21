package envfile

import (
	"fmt"
	"os"

	envv1 "go.firedancer.io/radiance/proto/env/v1"
	"google.golang.org/protobuf/encoding/prototext"
)

// Load loads the environment config from the given prototxt file.
func Load(filename string) (*envv1.Env, error) {
	b, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read env file: %w", err)
	}

	var env envv1.Env
	err = prototext.Unmarshal(b, &env)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal env file: %w", err)
	}

	return &env, nil
}
