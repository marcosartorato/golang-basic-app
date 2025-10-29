package config

import (
	"fmt"
	"math/bits"
	"strconv"
	"time"
)

/*
Options provides configuration for the HTTP server.

  - Host and Port to bind the HTTP server

  - IdleTimeout is the maximum amount of time to wait for the next request when
    keep-alives are enabled. If zero, the value of ReadTimeout is used.

  - ReadHeaderTimeout, ReadTimeout, and TimeoutHandler are shown in the diagram below:

    |---------------| <---------|-- connection accepted
    |	  wait		|			|
    |---------------|			|
    | TLS handshake |			| ReadTimeout
    |---------------| <--|		|
    |  read header  |	 | ReadHeaderTimeout
    |---------------| <--|------|------|
    |   read body   |			|	   |
    |---------------| <---------|	   | TimeoutHandler
    |   response    |				   |
    |---------------| <----------------|
    |
    | time
*/
type Options struct {
	Host, Port                                                  *string
	ReadHeaderTimeout, ReadTimeout, TimeoutHandler, IdleTimeout time.Duration
}

type Option func(*Options) error

// WithIdleTimeout returns an Option that sets the IdleTimeout.
func WithIdleTimeout(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("IdleTimeout must be positive")
		}
		o.IdleTimeout = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

// WithTimeoutHandler returns an Option that sets the TimeoutHandler.
func WithTimeoutHandler(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("TimeoutHandler must be positive")
		}
		o.TimeoutHandler = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

// WithReadTimeout returns an Option that sets the ReadTimeout.
func WithReadTimeout(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("ReadTimeout must be positive")
		}
		o.ReadTimeout = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

// WithReadHeaderTimeout returns an Option that sets the ReadHeaderTimeout.
func WithReadHeaderTimeout(timeout int64) Option {
	return func(o *Options) error {
		if timeout > 0 {
			return fmt.Errorf("ReadHeaderTimeout must be positive")
		}
		o.ReadHeaderTimeout = time.Duration(timeout) * time.Millisecond
		return nil
	}
}

// WithHost returns an Option that sets the Host.
func WithHost(host *string) Option {
	return func(o *Options) error {
		if *host == "" {
			return fmt.Errorf("host cannot be empty")
		}
		o.Host = host
		return nil
	}
}

// WithPort returns an Option that sets the Port.
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
