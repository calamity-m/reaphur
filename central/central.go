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
	"github.com/calamity-m/reaphur/central/internal/srv"
	"github.com/calamity-m/reaphur/central/internal/util"
	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/calamity-m/reaphur/pkg/logging"
	"github.com/spf13/cobra"
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
				os.Exit(1)
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

			oa := util.CreateNewOpenAIClient(cfg.AIToken)

			server, err := srv.NewCentralServiceServer(
				logger,
				cfg,
				parser.NewOpenAIParser(logger, oa),
				fncall.NewOpenAIFnCaller(logger, oa),
				persistence.NewMemoryFoodStore(),
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
)
