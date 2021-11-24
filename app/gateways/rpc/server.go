package rpc

import (
	"context"
	"fmt"
	"log"
	"net/http"

	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpcRecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/newrelic/go-agent/v3/integrations/nrgrpc"
	"github.com/newrelic/go-agent/v3/newrelic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stone-co/the-amazing-ledger/app"
	"github.com/stone-co/the-amazing-ledger/app/domain/usecases"
	httpHandlers "github.com/stone-co/the-amazing-ledger/app/gateways/http"
	proto "github.com/stone-co/the-amazing-ledger/gen/ledger/v1beta"
)

func NewServer(ctx context.Context, useCase *usecases.LedgerUseCase, nr *newrelic.Application, cfg *app.Config, commit, time string) (*grpc.Server, *http.Server, error) {
	api := NewAPI(useCase)

	grpcServer := newRPCServer(api, nr)

	server, err := newGatewayServer(ctx, cfg, commit, time)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create new GRPC server: %w", err)
	}

	return grpcServer, server, nil
}

func newRPCServer(api *API, nr *newrelic.Application) *grpc.Server {
	// Define a func to handle panic
	dealPanic := func(p interface{}) (err error) {
		log.Printf("panic triggered: %v", p)
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opts := []grpcRecovery.Option{
		grpcRecovery.WithRecoveryHandler(dealPanic),
	}

	srv := grpc.NewServer(
		grpcMiddleware.WithUnaryServerChain(
			grpcRecovery.UnaryServerInterceptor(opts...),
			nrgrpc.UnaryServerInterceptor(nr),
			loggerInterceptor,
		),
		grpcMiddleware.WithStreamServerChain(
			grpcRecovery.StreamServerInterceptor(opts...),
			nrgrpc.StreamServerInterceptor(nr),
		),
	)

	proto.RegisterLedgerAPIServer(srv, api)
	proto.RegisterHealthAPIServer(srv, api)

	return srv
}

func newGatewayServer(ctx context.Context, cfg *app.Config, commit, time string) (*http.Server, error) {
	gwMux := runtime.NewServeMux()
	gwEndpoint := fmt.Sprintf("%s:%d", cfg.RPCServer.Host, cfg.RPCServer.Port)

	conn, err := grpc.DialContext(ctx, gwEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, fmt.Errorf("failed to dial server: %w", err)
	}

	err = proto.RegisterLedgerAPIHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register ledger handler: %w", err)
	}

	err = proto.RegisterHealthAPIHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, fmt.Errorf("failed to register health handler: %w", err)
	}

	err = gwMux.HandlePath(http.MethodGet, "/metrics", httpHandlers.MetricsHandler)
	if err != nil {
		return nil, fmt.Errorf("failed to configure metrics handler: %w", err)
	}

	err = gwMux.HandlePath(http.MethodGet, "/version", httpHandlers.VersionHandler(commit, time))
	if err != nil {
		return nil, fmt.Errorf("failed to configure version handler: %w", err)
	}

	gwServer := &http.Server{
		Addr:         fmt.Sprintf("%s:%d", cfg.HttpServer.Host, cfg.HttpServer.Port),
		Handler:      gwMux,
		ReadTimeout:  cfg.HttpServer.ReadTimeout,
		WriteTimeout: cfg.HttpServer.WriteTimeout,
	}

	return gwServer, nil
}
