package frontdoor

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

type Site struct {
	r *mux.Router
}

func NewSite(logger logrus.FieldLogger) *Site {
	r := mux.NewRouter()

	homeHandler := wrapWithLogger("frontdoor", "home", newHomeHandler(logger), logger)
	r.Handle("/", homeHandler).Methods(http.MethodGet)

	updateHandler := wrapWithLogger("frontdoor", "update", newUpdateHandler(logger), logger)
	r.Handle("/update", updateHandler).Methods(http.MethodPost)

	return &Site{r: r}
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
