package frontdoor

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type updateHandler struct {
	logger logrus.FieldLogger
}

func newUpdateHandler(logger logrus.FieldLogger) *updateHandler {
	if logger == nil {
		logger = logrus.New()
	}

	return &updateHandler{logger: logger}
}

func (h *updateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		h.logger.WithError(err).Error("parsing form")
		w.WriteHeader(http.StatusInternalServerError)

		_, err = w.Write([]byte("internal server error"))
		if err != nil {
			h.logger.WithError(err).Error("write error")
		}

		return
	}

	name := r.Form.Get("name")
	h.logger.WithField("name", name).Info("adding name to guestbook")

	http.Redirect(w, r, "/", http.StatusFound)
}
