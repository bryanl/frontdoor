package frontdoor

import (
	"net/http"
	"time"

	"github.com/sirupsen/logrus"
)

type LoggerMiddleware struct {
	// Logger is the logger instance used to log messages
	Logger logrus.FieldLogger
	// Name is the name of the application
	Name string
}

func (m *LoggerMiddleware) Handler(h http.Handler, component string) *LoggerHandler {
	return &LoggerHandler{
		m:         m,
		handler:   h,
		component: component,
	}
}

type responseData struct {
	status int
	size   int
}

type LoggerHandler struct {
	http.ResponseWriter
	m            *LoggerMiddleware
	handler      http.Handler
	component    string
	responseData *responseData
}

func (h *LoggerHandler) newResponseData() *responseData {
	return &responseData{}
}

func (h *LoggerHandler) Write(data []byte) (int, error) {
	if h.responseData.status == 0 {
		// ensure we have a return status
		h.responseData.status = http.StatusOK
	}

	size, err := h.ResponseWriter.Write(data)
	h.responseData.size += size
	return size, err
}

func (h *LoggerHandler) WriteHeader(s int) {
	h.ResponseWriter.WriteHeader(s)
	h.responseData.status = s
}

func (h *LoggerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	h = h.m.Handler(h.handler, h.component)
	h.ResponseWriter = w
	h.responseData = h.newResponseData()
	h.handler.ServeHTTP(h, r)

	elapsed := time.Since(start)

	status := h.responseData.status
	if status == 0 {
		status = http.StatusOK
	}

	fields := logrus.Fields{
		"status":     status,
		"method":     r.Method,
		"request":    r.RequestURI,
		"remote":     r.RemoteAddr,
		"duration":   float64(elapsed.Nanoseconds()) / float64(1000),
		"size":       h.responseData.size,
		"referer":    r.Referer(),
		"user-agent": r.UserAgent(),
	}

	if h.m.Name != "" {
		fields["name"] = h.m.Name
	}

	if h.component != "" {
		fields["component"] = h.component
	}

	if l := h.m.Logger; l != nil {
		l.WithFields(fields).Info("completed handling request")
	} else {
		logrus.WithFields(fields).Info("completed handling request")
	}
}
