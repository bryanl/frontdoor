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

	h := &site{}
	server := newServer(logger, addr, h)

	go func(server *http.Server) {
		logger.WithFields(logrus.Fields{
			"addr": server.Addr,
		}).Info("service is listening")

		if err := server.ListenAndServe(); err != http.ErrServerClosed {
			logrus.WithError(err).Fatal("service received an unexpected error")
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

func newServer(logger logrus.FieldLogger, addr string, h http.Handler) *http.Server {
	if addr == "" {
		addr = defaultAddr
	}

	m := frontdoor.LoggerMiddleware{
		Logger: logger,
		Name:   "frontdoor",
	}

	mh := m.Handler(h, "home")

	return &http.Server{Addr: addr, Handler: mh}
}

type site struct{}

func (s *site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte(tmpl))
}

var tmpl = `
<html>
<head>Front Door</head>
<body>
<p>
Some down this will be a guestbook
</p>
</body>
</html>
`
