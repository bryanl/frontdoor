package frontdoor

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type updateHandler struct {
	logger     logrus.FieldLogger
	repository Repository
}

func newUpdateHandler(logger logrus.FieldLogger, r Repository) *updateHandler {
	if logger == nil {
		logger = logrus.New()
	}

	return &updateHandler{
		logger:     logger,
		repository: r,
	}
}

func (h *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.WithError(err).Error("parsing form")
		http.Error(w, "parsing form", http.StatusInternalServerError)
		return
	}

	name := r.Form.Get("name")
	if err := h.repository.AddName(r.Context(), name); err != nil {
		h.logger.WithError(err).
			WithField("name", name).Error("adding name to list")
		http.Error(w, "unable to add name", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusFound)
}
