package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"mime"
	"net/http"
	"strings"
	internalgrpc "triton-grpc-gateway/internal/grpc"
	tls2 "triton-grpc-gateway/pkg/tls"

	_ "triton-grpc-gateway/statik"
)

func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")

	statikFS, err := fs.New()
	if err != nil {
		// Panic since this is a permanent error.
		panic("creating OpenAPI filesystem: " + err.Error())
	}

	return http.FileServer(statikFS)
}

func run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Register gRPC server endpoint
	gwmux := runtime.NewServeMux()

	tlsConfig, err := tls2.NewConfigTLS("./conf/tls-self.crt", "./conf/tls-self.key").BuildTLSConfig(tls2.ScalaServiceCipherSuites())
	if err != nil {
		logrus.Fatalf("build tls config failed: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
	}

	// TODO: correct endpoint
	err = internalgrpc.RegisterUAAServiceHandlerFromEndpoint(ctx, gwmux, "uaa.triton-transparency.com:443", opts)
	// err = internalgrpc.RegisterUAAServiceHandlerFromEndpoint(ctx, gwmux, "localhost:8081", opts)
	if err != nil {
		return err
	}

	// swagger endpoint
	oa := getOpenAPIHandler()

	gatewayAddr := fmt.Sprintf("%s:%d", "", 8080)

	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, "/api") {
				gwmux.ServeHTTP(w, r)
				return
			}

			oa.ServeHTTP(w, r)
		}),
	}

	return fmt.Errorf("serving gRPC-Gateway server: %w", gwServer.ListenAndServe())
}

func main() {
	flag.Parse()

	logrus.SetLevel(logrus.DebugLevel)

	if err := run(); err != nil {
		logrus.Fatal(err)
	}
}
