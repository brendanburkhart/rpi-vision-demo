package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	api "github.com/brendanburkhart/rpi-vision-demo/webserver/internal/api"
	fileserver "github.com/brendanburkhart/rpi-vision-demo/webserver/internal/fileserver"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		for range c {
			cancel()
		}
	}()

	if err := run(ctx); err != nil {
		fmt.Printf("got error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	api := &api.Server{
		ListenAddress: ":8181",
		GrpcAddress:   "localhost:50051",
		Version:       "0.1.0",
	}

	fs := &fileserver.Server{
		ListenAddress: ":8080",
		Dir:           "./dashboard",
	}

	serverCtx, serverCancel := context.WithCancel(ctx)
	defer func() {
		serverCancel()
	}()

	errs := make(chan error)
	go func() {
		errs <- fs.Run(serverCtx)
	}()

	go func() {
		errs <- api.Run(serverCtx)
	}()

	return <-errs
}
