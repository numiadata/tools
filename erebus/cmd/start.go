package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/neilotoole/errgroup"
	"github.com/numiadata/tools/erebus/config"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start erebus state streaming file consumer",
		Long: `Start erebus state streaming file consumer. The process will watch
the configured state streaming directory for changes. New files will be parsed
and sent to the configured sink. Existing files may have been ignored, so erebus
will attempt to examine missing consumed files and handle them appropriately.
`,
		RunE: startCmdHandler,
	}

	_ = cmd.MarkFlagRequired(flagConfig)
	_ = cmd.MarkFlagRequired(flagStateStreamingDir)

	return cmd
}

func startCmdHandler(cmd *cobra.Command, args []string) error {
	cfg, err := config.Parse(cfgFile)
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	if _, err := os.Stat(stateStreamDir); os.IsNotExist(err) {
		return fmt.Errorf("state streaming directory '%s' does not exist", stateStreamDir)
	}

	logger, err := getCmdLogger(cmd)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithCancel(cmd.Context())
	g, ctx := errgroup.WithContext(ctx)

	// listen for and trap any OS signal to gracefully shutdown and exit
	trapSignal(cancel, logger)

	g.Go(func() error {
		return watchStreamingDir(ctx, logger)
	})

	// Block main process until all spawned goroutines have gracefully exited and
	// signal has been captured in the main process or if an error occurs.
	return g.Wait()
}

func createFileWatcher() (*watcher.Watcher, error) {
	w := watcher.New()

	// SetMaxEvents to 1 to allow at most 1 event's to be received
	// on the Event channel per watching cycle.
	//
	// If SetMaxEvents is not set, the default is to send all events.
	w.SetMaxEvents(1)

	// Only notify when files are written to, this will allow us to capture new
	// files and files currently being written to.
	w.FilterOps(watcher.Write)

	// Only files that match the regular expression during file listings will be
	// watched.
	r := regexp.MustCompile(fmt.Sprintf("^%s", ssFilePrefix))
	w.AddFilterHook(watcher.RegexFilterHook(r, false))

	if err := w.Add(stateStreamDir); err != nil {
		return nil, fmt.Errorf("failed to add %s to directory watcher: %v", stateStreamDir, err)
	}

	return w, nil
}

func watchStreamingDir(ctx context.Context, logger zerolog.Logger) error {
	w, err := createFileWatcher()
	if err != nil {
		return err
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- w.Start(time.Millisecond * 100)
	}()

	logger.Info().Str("dir", stateStreamDir).Str("file_prefix", ssFilePrefix).Msg("watching state streaming directory")

	for {
		select {
		case event := <-w.Event:
			fmt.Printf("%+v\n", event)

		case <-ctx.Done():
			// Context was explicitly cancelled due to a signal capture so we can safely
			// close the the watcher. This will cause the watch process to safely exit.
			w.Close()
			return ctx.Err()

		case err := <-w.Error:
			logger.Error().Err(err).Msg("directory watch failure")
			return err

		case <-w.Closed:
			return nil

		case err := <-errCh:
			logger.Error().Err(err).Msg("directory watch failure")
			return err
		}
	}
}
