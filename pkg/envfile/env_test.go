package envfile

import (
	"testing"

	envv1 "go.firedancer.io/radiance/proto/env/v1"
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

	if len(env.Kafka.Brokers) != 2 {
		t.Errorf("Expected 2 broker, got %d", len(env.Kafka.Brokers))
	}

	if _, ok := env.Kafka.Encryption.(*envv1.Kafka_TlsEncryption); !ok {
		t.Errorf("Expected TLS encryption, got %T", env.Kafka.Encryption)
	}

}
