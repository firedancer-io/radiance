package envfile

import (
	"testing"
)

func TestLoadEnvFile(t *testing.T) {
	env, err := Load("testdata/env.prototxt")
	if err != nil {
		t.Fatal(err)
	}

	if env.Nodes[0].Name != "localhost" {
		t.Errorf("Expected node name to be 'localhost', got '%s'", env.Nodes[0].Name)
	}

	if len(env.Nodes) != 2 {
		t.Errorf("Expected 2 node, got %d", len(env.Nodes))
	}
}
