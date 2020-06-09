package server

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"time"

	pb "github.com/brendanburkhart/rpi-vision-demo/webserver/grpc-common"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) getThresholdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		threshold, err := s.grpcClient.GetThresholds(ctx, &emptypb.Empty{})
		if err != nil {
			respondError(w, http.StatusInternalServerError)
			log.Printf("could not get threshold: %v", err)
		}

		respond(w, &threshold, http.StatusOK)
	}
}

func (s *Server) postThresholdHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var threshold pb.Thresholds
		if err := json.NewDecoder(r.Body).Decode(&threshold); err != nil {
			respondError(w, http.StatusUnprocessableEntity)
			return
		}

		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err := s.grpcClient.SetThresholds(ctx, &threshold)
		if err != nil {
			respondError(w, http.StatusInternalServerError)
			log.Printf("could not set threshold: %v", err)
		}

		w.WriteHeader(http.StatusOK)
	}
}
