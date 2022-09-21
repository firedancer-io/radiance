package kafka

import (
	"crypto/tls"
	"net"
	"time"

	envv1 "go.firedancer.io/radiance/proto/env/v1"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/scram"
)

func NewClientFromEnv(env *envv1.Kafka, opts ...kgo.Opt) (*kgo.Client, error) {
	opts = append(opts,
		kgo.SeedBrokers(env.Brokers...),
	)

	if s, ok := env.Auth.(*envv1.Kafka_SaslAuth); ok {
		m := scram.Auth{
			User: s.SaslAuth.Username,
			Pass: s.SaslAuth.Password}.AsSha256Mechanism()

		opts = append(opts, kgo.SASL(m))
	}

	if _, ok := env.Encryption.(*envv1.Kafka_TlsEncryption); ok {
		tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		opts = append(opts, kgo.Dialer(tlsDialer.DialContext))
	}

	cl, err := kgo.NewClient(opts...)
	if err != nil {
		return nil, err
	}

	return cl, nil
}
