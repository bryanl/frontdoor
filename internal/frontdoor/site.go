package frontdoor

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type Site struct {
	r *mux.Router
}

func NewSite(logger logrus.FieldLogger, repo Repository) (*Site, error) {
	r := mux.NewRouter()

	homeHandler, err := newHomeHandler(logger, repo)
	if err != nil {
		return nil, errors.Wrap(err, "initializing home handler")
	}

	wrapWithLogger("frontdoor", "home", homeHandler, logger)
	r.Handle("/", homeHandler).Methods(http.MethodGet)

	updateHandler := wrapWithLogger("frontdoor", "update", newUpdateHandler(logger, repo), logger)
	r.Handle("/update", updateHandler).Methods(http.MethodPost)

	healthzHandler := wrapWithLogger("frontdoor", "healthz", newHealthzHandler(repo), logger)
	r.Handle("/healthz", healthzHandler).Methods(http.MethodGet)

	return &Site{r: r}, nil
}

func (s *Site) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}

func wrapWithLogger(name, componentName string, h http.Handler, logger logrus.FieldLogger) http.Handler {
	m := LoggerMiddleware{
		Logger: logger,
		Name:   name,
	}

	return m.Handler(h, componentName)
}
