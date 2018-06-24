package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/bryanl/frontdoor/internal/frontdoor"

	"github.com/sirupsen/logrus"
)

const (
	defaultAddr = ":8080"
)

func main() {
	logger := logrus.New()

	addr := os.Getenv("FRONTDOOR_ADDR")
	redisAddr := os.Getenv("FRONTDOOR_REDIS_ADDR")

	repo := frontdoor.NewRedisRepository(redisAddr, logger)

	h, err := frontdoor.NewSite(logger, repo)
	if err != nil {
		logger.WithError(err).Error("initializing site")
		os.Exit(1)
	}

	server := newServer(addr, h)

	go func(server *http.Server) {
		logger.WithFields(logrus.Fields{
			"addr": server.Addr,
		}).Info("service is listening")

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logrus.WithError(err).Error("service received an unexpected error")
			os.Exit(1)
		}
	}(server)

	graceful(server, logger, 5*time.Second)
}

func graceful(server *http.Server, logger logrus.FieldLogger, timeout time.Duration) {
	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	logger.WithField("timeout", timeout).Info("shutting down")

	if err := server.Shutdown(ctx); err != nil {
		logger.WithError(err).Error("unable to shutdown")
		return
	}

	logger.Info("service stopped")
}

func newServer(addr string, h http.Handler) *http.Server {
	if addr == "" {
		addr = defaultAddr
	}

	return &http.Server{Addr: addr, Handler: h}
}
