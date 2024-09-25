package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
	"os/signal"
	"time"
)

const (
	defaultReadTimeout  = 30 * time.Second
	defaultWriteTimeout = 30 * time.Second
	defaultPort         = "8080"
)

type ShutdownFunc func(ctx context.Context) error

type Option func(config *serverConfig)

type serverConfig struct {
	port         string
	readTimeout  time.Duration
	writeTimeout time.Duration
	shutdownFunc []ShutdownFunc
}

// SetPort - sets the port for the server.  Defaults to 8080 if not set
func SetPort(port string) Option {
	return func(c *serverConfig) {
		c.port = port
	}
}

// SetReadTimeout - sets the read timeout for the server.  Defaults to 30s if not set
func SetReadTimeout(to time.Duration) Option {
	return func(c *serverConfig) {
		c.readTimeout = to
	}
}

// SetWriteTimeout - sets the write timeout for the server.  Defaults to 30s if not set
func SetWriteTimeout(to time.Duration) Option {
	return func(c *serverConfig) {
		c.writeTimeout = to
	}
}

// SetShutdownFuncs - sets a number of ShutdownFunc to be executed before the
// server comes down.  No defaults
func SetShutdownFuncs(fns ...ShutdownFunc) Option {
	return func(c *serverConfig) {
		c.shutdownFunc = fns
	}
}

// StartServer - creates a http server with the incoming configuration.  For any
// unset config, a default will be used.  This method
// will catch an interrupt even and execute the list of shutdown functions provided
// before exiting.  These functions can be any clean up methods which should run
// before the entire service shuts down.  This is a blocking method.
func StartServer(ctx context.Context, r *mux.Router, opts ...Option) error {

	// set the defaults
	config := &serverConfig{
		port:         defaultPort,
		readTimeout:  defaultReadTimeout,
		writeTimeout: defaultWriteTimeout,
		shutdownFunc: nil,
	}

	for _, opt := range opts {
		opt(config)
	}

	var errs error
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)

	defer stop()

	defer func(c context.Context) {
		if config.shutdownFunc != nil {
			for _, f := range config.shutdownFunc {
				if err := f(c); err != nil {
					errs = errors.Join(errs, err)
					logrus.Warnf("error thrown from shutdown fn %v: %v", f, err)
				}
			}
		}
	}(ctx)

	serverPort := fmt.Sprintf(":%s", config.port)

	logrus.Infof("starting http server on port %s", serverPort)

	srv := &http.Server{
		Addr:         serverPort,
		Handler:      r,
		ReadTimeout:  config.readTimeout,
		WriteTimeout: config.writeTimeout,
	}

	srvErr := make(chan error, 1)

	go func() {
		srvErr <- srv.ListenAndServe()
	}()

	// wait here until the server errors or comes down
	select {
	case errs = <-srvErr:
		return errs
	case <-ctx.Done():
		stop()
	}

	if e := srv.Shutdown(ctx); e != nil {
		errs = errors.Join(errs, e)
	}

	return errs
}
