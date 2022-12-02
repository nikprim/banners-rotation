package cmd

import (
	"context"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var (
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "banner-rotator",
		Short: "banner-rotator app",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Usage()
		},
	}
)

func Execute(ctx context.Context) {
	rootCmd.AddCommand(serveHTTPCommand(ctx))

	if err := rootCmd.Execute(); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
