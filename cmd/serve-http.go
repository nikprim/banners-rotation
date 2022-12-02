package cmd

import (
	"context"
	"errors"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nikprim/banners-rotation/cmd/config"
	"github.com/nikprim/banners-rotation/internal/app"
	"github.com/nikprim/banners-rotation/internal/rmq"
	internalGRPC "github.com/nikprim/banners-rotation/internal/server/grpc"
	psqlStorage "github.com/nikprim/banners-rotation/internal/storage"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

func serveHTTPCommand(ctx context.Context) *cobra.Command {
	command := &cobra.Command{
		Use:   "serve-http",
		Short: "serves http api",
		RunE:  serveHTTPCommandRunE(ctx),
	}

	command.Flags().StringVar(&cfgFile, "config", "", "Path to configuration file")

	err := command.MarkFlagRequired("config")
	if err != nil {
		return nil
	}

	return command
}

func serveHTTPCommandRunE(ctx context.Context) func(cmd *cobra.Command, args []string) (err error) {
	return func(cmd *cobra.Command, args []string) (err error) {
		configFile := cmd.Flag("config").Value.String()

		cfg, err := config.ParseBannerRotatorConfig(configFile)
		if err != nil {
			log.Error().Err(err).Msg("failed to parse config")

			return err
		}

		logLevel, err := zerolog.ParseLevel(cfg.Logger.Level)
		if err != nil {
			log.Error().Err(err).Msg("failed to install log level")

			return err
		}

		zerolog.SetGlobalLevel(logLevel)

		ctx, cancel := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		defer cancel()

		conn, err := pgxpool.Connect(ctx, cfg.DB.URI)
		if err != nil {
			log.Error().Err(err).Msg("unable to connect to database")

			return err
		}

		store := psqlStorage.New(conn)

		cfgP := cfg.Producer
		producer := rmq.NewProducer(cfgP.URI, cfgP.Queue)

		err = producer.Connect()
		if err != nil {
			log.Error().Err(err).Msg("unable to connect to amqp")

			return err
		}

		application := app.New(store, producer)
		grpcServer := internalGRPC.NewServer(cfg.GRPC.Host, cfg.GRPC.Port, application)

		go func() {
			<-ctx.Done()

			log.Info().Msg("disconnecting a db...")

			ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			conn.Close()

			log.Info().Msg("db is disconnected")
			log.Info().Msg("disconnecting a producer...")

			ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			if err := producer.Disconnect(); err != nil {
				log.Error().Err(err).Msg("failed to disconnect producer")
			}

			log.Info().Msg("producer is disconnected")
			log.Info().Msg("stopping an grpc server...")

			ctx, cancel = context.WithTimeout(context.Background(), time.Second*2)
			defer cancel()

			if err := grpcServer.Stop(ctx); err != nil {
				log.Error().Err(err).Msg("failed to stop grpc server")
			}

			log.Info().Msg("grpc server is stopped")
		}()

		log.Info().Msg("banner-rotator is running...")

		if err := grpcServer.Start(ctx); err != nil {
			log.Info().Msg("123...")
			cancel()

			if !errors.Is(err, http.ErrServerClosed) {
				log.Error().Err(err).Msg("failed to start grpc server")
			}
		}

		return nil
	}
}
