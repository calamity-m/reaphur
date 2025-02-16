package srv

import (
	"fmt"
	"log/slog"
	"net"
	"os"

	"github.com/calamity-m/reaphur/central/internal/conf"
	"github.com/calamity-m/reaphur/central/internal/fncall"
	"github.com/calamity-m/reaphur/central/internal/parser"
	"github.com/calamity-m/reaphur/central/internal/persistence"
	"github.com/calamity-m/reaphur/pkg/errs"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type CentralServiceServer struct {
	logger *slog.Logger

	config *conf.Config

	parser *parser.OpenAIParser

	fnCaller *fncall.OpenAIFnCaller

	foodStore persistence.FoodPersistence

	centralproto.UnimplementedCentralServiceServer
	centralproto.UnimplementedCentralFoodServiceServer
}

// Runs the GRPC server until notify is pushed to. You can wait
// on the returned channel for an exit code to account for graceful
// shutdowns.
func (s *CentralServiceServer) Run(notify <-chan os.Signal) <-chan error {
	// Channel we'll use to signal for finish
	exit := make(chan error)

	// create the grpc server that we can later serve on
	grpcServer := grpc.NewServer(s.config.GrpcServerOpts...)

	// register ourselves
	centralproto.RegisterCentralServiceServer(grpcServer, s)
	centralproto.RegisterCentralFoodServiceServer(grpcServer, s)

	if s.config.Reflect {
		reflection.Register(grpcServer)
	}

	go func() {
		// At the end of our function. If no errors were otherwise
		// pushed to this channel, it notifies as a successful shutdown
		defer close(exit)

		// Block until we receive a Interrupt or Kill
		<-notify

		grpcServer.GracefulStop()
	}()

	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.logger.Error(fmt.Sprintf("failed to listen on: %q", s.config.Address), slog.Any("err", err))
	}

	s.logger.Info(
		"Server has follwing configuration",
		slog.String("environment", s.config.Environment),
		slog.Bool("reflection", s.config.Reflect),
		slog.String("log_level", s.config.LogLevel.String()),
		slog.Bool("log_structured", s.config.LogStructured),
		slog.Bool("log_request_id", s.config.LogRequestId),
	)
	s.logger.Info(fmt.Sprintf("Starting server on %s", s.config.Address))
	if err := grpcServer.Serve(listener); err != nil {
		exit <- fmt.Errorf("failed to start/close sow server due to: %w", err)
	}

	return exit
}

func NewCentralServiceServer(logger *slog.Logger, config *conf.Config, openai *parser.OpenAIParser, fnCaller *fncall.OpenAIFnCaller, foodStore persistence.FoodPersistence) (*CentralServiceServer, error) {
	if logger == nil || config == nil {
		return nil, errs.ErrNilNotAllowed
	}

	s := &CentralServiceServer{
		logger:    logger,
		config:    config,
		foodStore: foodStore,
		parser:    openai,
		fnCaller:  fnCaller,
	}

	return s, nil
}

func (s *CentralServiceServer) commonServiceValidation() error {
	if s.logger == nil {
		return errs.ErrNilNotAllowed
	}
	if s.config == nil {
		return errs.ErrNilNotAllowed
	}
	if s.parser == nil {
		return errs.ErrNilNotAllowed
	}
	if s.fnCaller == nil {
		return errs.ErrNilNotAllowed
	}

	return nil
}
