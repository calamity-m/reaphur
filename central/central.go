package central

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/calamity-m/reaphur/central/internal/conf"
	"github.com/calamity-m/reaphur/central/internal/fncall"
	"github.com/calamity-m/reaphur/central/internal/parser"
	"github.com/calamity-m/reaphur/central/internal/persistence"
	"github.com/calamity-m/reaphur/central/internal/prompts"
	"github.com/calamity-m/reaphur/central/internal/srv"
	"github.com/calamity-m/reaphur/central/internal/util"
	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/calamity-m/reaphur/pkg/logging"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

var (
	CentralCommand = &cobra.Command{
		Use:   "central",
		Short: "central short",
		Long:  `central long`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := conf.NewConfig(bindings.Debug)
			if err != nil {
				fmt.Printf("Failed to create config: %v\n", err)
				return err
			}

			// Create logger
			logger := slog.New(logging.NewCustomizedHandler(os.Stderr, &logging.CustomHandlerCfg{
				Structed:        cfg.LogStructured,
				RecordRequestId: cfg.LogRequestId,
				Level:           cfg.LogLevel,
				AddSource:       cfg.LogAddSource,
				StaticAttributes: []slog.Attr{
					slog.String("system", "reap"),
					slog.String("environment", cfg.Environment),
				},
			}))

			// Display some helpful starting info
			logger.Info(fmt.Sprintf("Redis Address: %s", cfg.RedisAddress))
			logger.Info(fmt.Sprintf("GRPC Reflection: %t", cfg.Reflect))
			logger.Info(fmt.Sprintf("Environment: %s", cfg.Environment))

			// Display schemas
			logger.Debug(
				"scheams",
				slog.String("create_food", prompts.CreateCardioJson),
				slog.String("create_weight_lifting", prompts.CreateWeightLiftingJson),
				slog.String("create_cardio", prompts.CreateCardioJson),
				slog.String("get_food", prompts.GetFoodJson),
			)

			oa := util.CreateNewOpenAIClient(cfg.AIToken)
			foodStore, err := persistence.NewRedisFoodStore(logger, cfg)
			if err != nil {
				logger.Error("failed to create redis food store", slog.Any("err", err))
				return err
			}

			server, err := srv.NewCentralServiceServer(
				logger,
				cfg,
				parser.NewOpenAIParser(logger, oa),
				fncall.NewOpenAIFnCaller(logger, oa),
				foodStore,
			)
			if err != nil {
				logger.Error("failed to run server", slog.Any("err", err))
				return err
			}

			// Create the channel which will wait for our shutdown signals which we can
			// utilise for a nice graceful shutdown
			sig := make(chan os.Signal, 2)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			//  Run the server
			exit := server.Run(sig)

			// Check the exit error
			if err = <-exit; err != nil {
				logger.Error("Received an error on server exit", slog.Any("err", err))
				return err
			}

			return nil
		},
	}

	CentralGenerateSchemaCommand = &cobra.Command{
		Use:   "generate",
		Short: "generate schemas",
		Long:  `generate function schemas for use with smarty pants services`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return prompts.GenerateSchemas()
		},
	}
)

func NewCentralServiceClient(addr string, opts []grpc.DialOption) (centralproto.CentralServiceClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(addr, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to dial: %v", err)
	}
	client := centralproto.NewCentralServiceClient(conn)
	return client, conn, nil
}
