package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	pb "github.com/brendanburkhart/rpi-vision-demo/webserver/grpc-common"
	"google.golang.org/grpc"
)

// Server is the vision API server
type Server struct {
	ListenAddress string
	GrpcAddress   string
	Version       string
	start         time.Time
	grpcClient    pb.PipelineControllerClient
}

func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 1e3)
		next.ServeHTTP(w, r)
	})
}

// Run web server
func (s *Server) Run(ctx context.Context) error {
	router := s.registerRoutes()

	var handler http.Handler = router
	handler = limitBody(handler)

	httpServer := &http.Server{
		Addr:              s.ListenAddress,
		Handler:           handler,
		ReadTimeout:       time.Second * 15,
		ReadHeaderTimeout: time.Second * 15,
		WriteTimeout:      time.Second * 15,
		IdleTimeout:       time.Second * 30,
		MaxHeaderBytes:    4096,
	}

	grpcConn, err := grpc.Dial(s.GrpcAddress, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("did not connect: %v", err)
	}
	defer grpcConn.Close()

	s.grpcClient = pb.NewPipelineControllerClient(grpcConn)

	s.start = time.Now()

	errs := make(chan error)
	go func() {
		errs <- httpServer.ListenAndServe()
	}()

	select {
	case err := <-errs:
		return err
	case <-ctx.Done():
		return httpServer.Shutdown(ctx)
	}
}
