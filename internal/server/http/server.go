package internalhttp

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/julinserg/OtusAlgorithmHomeProject/internal/app"
)

type Application interface {
	GetAllDocument() ([]app.Document, error)
	AddNewDocument(url string) ([]app.Document, error)
	Search(str string) ([]app.SearchResult, error)
}

type Server struct {
	server   *http.Server
	logger   Logger
	endpoint string
}

type Logger interface {
	Info(msg string)
	Error(msg string)
	Debug(msg string)
	Warn(msg string)
}

type StatusRecorder struct {
	http.ResponseWriter
	Status int
}

func (r *StatusRecorder) WriteHeader(status int) {
	r.Status = status
	r.ResponseWriter.WriteHeader(status)
}

func NewServer(logger Logger, app Application, endpoint string) *Server {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:              endpoint,
		Handler:           loggingMiddleware(mux, logger),
		ReadHeaderTimeout: 3 * time.Second,
	}
	ch := minisearchHandler{logger: logger, app: app}
	mux.HandleFunc("/", ch.landingHandler)
	mux.HandleFunc("/search", ch.searchHandler)
	mux.HandleFunc("/add", ch.addHandler)
	return &Server{server, logger, endpoint}
}

func (s *Server) Start(ctx context.Context) error {
	s.logger.Info("http server started on " + s.endpoint)
	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}
	s.logger.Info("http server stopped")
	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
