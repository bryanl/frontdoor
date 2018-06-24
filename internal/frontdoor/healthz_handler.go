package frontdoor

import "net/http"

type healthzHandler struct {
	repository Repository
}

func newHealthzHandler(r Repository) *healthzHandler {
	return &healthzHandler{
		repository: r,
	}
}

func (h *healthzHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.repository.Ready(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
