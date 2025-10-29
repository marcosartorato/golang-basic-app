package config

import (
	"fmt"
	"math/bits"
	"strconv"
	"time"
)

type Options struct {
	Host, Port                                     *string
	ReadHeaderTimeout, ReadTimeout, TimeoutHandler time.Duration
}

type Option func(*Options) error

func WithTimeoutHandler(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("TimeoutHandler must be positive")
		}
		o.TimeoutHandler = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

func WithReadTimeout(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("ReadTimeout must be positive")
		}
		o.ReadTimeout = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

func WithReadHeaderTimeout(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("ReadHeaderTimeout must be positive")
		}
		o.ReadHeaderTimeout = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

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
