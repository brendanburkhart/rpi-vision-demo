package static

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"
)

// Server is the static dashboard file server
type Server struct {
	ListenAddress string
	Dir           string
}

func containsDotFile(name string) bool {
	parts := strings.Split(name, "/")

	for _, part := range parts {
		if strings.HasPrefix(part, ".") {
			return true
		}
	}

	return false
}

type customFile struct {
	http.File
}

// Custom file system, currently just hides dot files
type customFileSystem struct {
	http.FileSystem
}

func (fs customFileSystem) Open(name string) (http.File, error) {
	if containsDotFile(name) {
		return nil, os.ErrPermission
	}

	file, err := fs.FileSystem.Open(name)
	if err != nil {
		return nil, err
	}

	return customFile{file}, err
}

// Run serves the dashboard files
func (s *Server) Run(ctx context.Context) error {
	fs := customFileSystem{http.Dir(s.Dir)}

	handler := http.NewServeMux()
	handler.Handle("/", http.FileServer(fs))

	fileServer := &http.Server{
		Addr:              s.ListenAddress,
		Handler:           handler,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 30,
		MaxHeaderBytes:    4096,
	}

	errs := make(chan error)
	go func() {
		errs <- fileServer.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return fileServer.Shutdown(ctx)
	}
}
