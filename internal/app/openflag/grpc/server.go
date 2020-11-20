package grpc

import (
	"context"
	"errors"
	"net"

	"github.com/OpenFlag/OpenFlag/internal/app/openflag/engine"
	"github.com/OpenFlag/OpenFlag/internal/app/openflag/model"
	"github.com/OpenFlag/OpenFlag/pkg/evaluation"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// ErrInternalServerError represents an error that we return when we have an error in the evaluation server.
var ErrInternalServerError = errors.New("internal server error")

// Server represents a struct for the gRPC server.
type Server struct {
	gRPCServer *grpc.Server
}

// New creates a new gRPC server.
func New(evaluationEngine engine.Engine, entityRepo model.EntityRepo) *Server {
	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			return status.Error(codes.Unknown, ErrInternalServerError.Error())
		}),
	}

	gRPCServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			metrics.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
	)

	evaluation.RegisterEvaluationServer(gRPCServer, &evaluationServer{
		Engine:     evaluationEngine,
		EntityRepo: entityRepo,
	})

	// Registers the server reflection service on the given gRPC server.
	reflection.Register(gRPCServer)

	// Should be called after all services have been registered with the server
	grpc_prometheus.Register(gRPCServer)

	return &Server{
		gRPCServer: gRPCServer,
	}
}

// Start starts a new gRPC server.
func (s *Server) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		logrus.Fatalf("failed to listen: %v", err)
	}

	return s.gRPCServer.Serve(lis)
}

// Shutdown shutdowns the gRPC server.
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func(done chan struct{}) {
		if s.gRPCServer == nil {
			close(done)
		}

		s.gRPCServer.GracefulStop()
		close(done)
	}(done)

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-done:
		return nil
	}
}
