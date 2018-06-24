package frontdoor

import (
	"html/template"
	"net/http"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type homeHandler struct {
	logger     logrus.FieldLogger
	repository Repository
	template   *template.Template
}

func newHomeHandler(logger logrus.FieldLogger, r Repository) (*homeHandler, error) {
	if logger == nil {
		logger = logrus.New()
	}

	tmpl, err := template.New("home").Parse(homeTemplate)
	if err != nil {
		return nil, errors.Wrap(err, "parsing home template")
	}

	return &homeHandler{
		logger:     logger,
		repository: r,
		template:   tmpl,
	}, nil
}

func (h *homeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")

	names, err := h.repository.ListNames(r.Context())
	if err != nil {
		http.Error(w, "unable to list names", http.StatusInternalServerError)
	}

	data := map[string]interface{}{
		"Names": names,
	}

	err = h.template.Execute(w, data)
	if err != nil {
		h.logger.WithError(err).Error("executing template")
	}
}

var homeTemplate = `
<html>
<head >
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

{{if .Names -}}
<div>
<ul>
	{{range .Names -}}
	<li>{{.}}</li>
	{{end -}}
</ul>
</div>
{{end -}}

</body>
</html>
`
