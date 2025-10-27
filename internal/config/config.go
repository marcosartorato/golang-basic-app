package config

import (
	"fmt"
	"math/bits"
	"strconv"
)

type Options struct {
	Host, Port *string
}

type Option func(*Options) error

func WithHost(host *string) Option {
	return func(o *Options) error {
		if *host == "" {
			return fmt.Errorf("host cannot be empty")
		}
		o.Host = host
		return nil
	}
}

func WithPort(port *string) Option {
	return func(o *Options) error {
		const (
			minPort, maxPort int64 = 1, 65535
			base             int   = 10
		)
		if portNmb, err := strconv.ParseInt(*port, base, bits.UintSize); err != nil || portNmb > maxPort || portNmb < minPort {
			return fmt.Errorf("invalid port: %s", *port)
		}
		o.Port = port
		return nil
	}
}
