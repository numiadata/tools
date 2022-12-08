package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	flagConfig            = "config"
	flagStateStreamingDir = "state-streaming-dir"
	flagLogLevel          = "log-level"
	flagLogFormat         = "log-format"
	flagFilePrefix        = "file-prefix"

	logLevelJSON = "json"
	logLevelText = "text"

	defaultFilePrefix = "block-"
)

func Execute() error {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	rootCmd.AddCommand(
		NewStartCmd(),
	)

	return rootCmd.Execute()
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erebus",
		Short: "Listen for Cosmos SDK state streaming files and publish them to a sink",
	}

	cmd.PersistentFlags().String(flagConfig, "", "Path to the configuration file")
	cmd.PersistentFlags().String(flagStateStreamingDir, "", "Path to the state streaming directory")
	cmd.PersistentFlags().String(flagLogLevel, zerolog.InfoLevel.String(), "logging level")
	cmd.PersistentFlags().String(flagLogFormat, logLevelText, "logging format; must be either json or text")
	cmd.PersistentFlags().String(flagFilePrefix, defaultFilePrefix, "File prefix for state streamed files")

	return cmd
}

func getCmdLogger(cmd *cobra.Command) (zerolog.Logger, error) {
	logLvlStr, err := cmd.Flags().GetString(flagLogLevel)
	if err != nil {
		return zerolog.Nop(), err
	}

	logLvl, err := zerolog.ParseLevel(logLvlStr)
	if err != nil {
		return zerolog.Nop(), err
	}

	logFormatStr, err := cmd.Flags().GetString(flagLogFormat)
	if err != nil {
		return zerolog.Nop(), err
	}

	var logWriter io.Writer
	switch strings.ToLower(logFormatStr) {
	case logLevelJSON:
		logWriter = os.Stderr

	case logLevelText:
		logWriter = zerolog.ConsoleWriter{Out: os.Stderr}

	default:
		return zerolog.Nop(), fmt.Errorf("invalid logging format: %s", logFormatStr)
	}

	return zerolog.New(logWriter).Level(logLvl).With().Timestamp().Logger(), nil
}

// trapSignal will listen for any OS signal and invoke the CancelFunc allowing
// the main process to gracefully exit.
func trapSignal(cancel context.CancelFunc, logger zerolog.Logger) {
	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGTERM)
	signal.Notify(sigCh, syscall.SIGINT)

	go func() {
		sig := <-sigCh
		logger.Info().Str("signal", sig.String()).Msg("caught signal; shutting down...")
		cancel()
	}()
}
