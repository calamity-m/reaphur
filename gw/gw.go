package gw

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/calamity-m/reaphur/gw/internal/conf"
	"github.com/calamity-m/reaphur/pkg/bindings"
	"github.com/calamity-m/reaphur/pkg/logging"
	centralproto "github.com/calamity-m/reaphur/proto/v1/central"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	GRPCGatewayCommand = &cobra.Command{
		Use:   "gw",
		Short: "gw short",
		Long:  `gw long`,
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

			return run(logger, cfg)
		},
	}
)

func run(logger *slog.Logger, cfg *conf.Config) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ssmux := http.NewServeMux()

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := centralproto.RegisterCentralServiceHandlerFromEndpoint(ctx, mux, bindings.DefaultCentralAddress, opts)
	if err != nil {
		return err
	}
	err = centralproto.RegisterCentralFoodServiceHandlerFromEndpoint(ctx, mux, bindings.DefaultCentralAddress, opts)
	if err != nil {
		return err
	}

	// mount a path to expose the generated OpenAPI specification on disk
	ssmux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./proto/v1/central/central.swagger.json")
	})

	ssmux.HandleFunc("/swagger-ui/swagger-food.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./proto/v1/central/central_food.swagger.json")
	})

	// mount the Swagger UI that uses the OpenAPI specification path above
	ssmux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./gw/swagger"))))

	ssmux.Handle("/", mux)

	// Start HTTP server (and proxy calls to gRPC server endpoint)
	logger.Info(fmt.Sprintf("Listening on %s", cfg.Address))
	return http.ListenAndServe(cfg.Address, ssmux)
}
