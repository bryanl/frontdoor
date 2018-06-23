package frontdoor

import (
	"net/http"

	"github.com/sirupsen/logrus"
)

type homeHandler struct {
	logger logrus.FieldLogger
}

func newHomeHandler(logger logrus.FieldLogger) *homeHandler {
	if logger == nil {
		logger = logrus.New()
	}

	return &homeHandler{logger: logger}
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	_, err := w.Write([]byte(homeTemplate))
	if err != nil {
		h.logger.WithError(err).Error("write error")
	}
}

var homeTemplate = `
<html>
<head>
<title>Front Door</title>
</head>
<body>
<p>
<form action="/update" method="post">
<p>
Sign the guestbook: <input type="text" name="name">
</p>

<input type="submit">
</p>
</body>
</html>
`
