package discord

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/calamity-m/reaphur/central"
	"github.com/calamity-m/reaphur/discord/internal/bot"
	"github.com/calamity-m/reaphur/discord/internal/conf"
	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/calamity-m/reaphur/pkg/logging"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	DiscordBotCommand = &cobra.Command{
		Use:   "discord",
		Short: "discord short",
		Long:  `discord long`,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := conf.NewConfig(bindings.Debug)
			if err != nil {
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

			logger.Info("initialized logging")

			opts := []grpc.DialOption{
				// For now just use insecure
				grpc.WithTransportCredentials(insecure.NewCredentials()),
			}
			centralClient, centralConn, err := central.NewCentralServiceClient(cfg.CentralServerAddress, opts)
			if err != nil {
				logger.Error("failed to create central client", slog.Any("err", err))
				return err
			}
			defer centralConn.Close()

			discordBot, err := bot.NewDiscordBot(logger, cfg, centralClient)
			if err != nil {
				logger.Error("failed to create discord bot", slog.Any("err", err))
				return err
			}

			// Create the channel which will wait for our shutdown signals which we can
			// utilise for a nice graceful shutdown
			sig := make(chan os.Signal, 2)
			signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

			if err := discordBot.Run(sig); err != nil {
				logger.Error("failed to run discord bot", slog.Any("err", err))
				return err
			}

			return nil
		},
	}
)
