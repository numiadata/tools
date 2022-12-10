package cmd

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/neilotoole/errgroup"
	"github.com/numiadata/tools/erebus/codec"
	"github.com/numiadata/tools/erebus/config"
	"github.com/numiadata/tools/erebus/io"
	"github.com/radovskyb/watcher"
	"github.com/rs/zerolog"
	"github.com/spf13/cobra"
)

const (
	watcherDuration     = time.Millisecond * 100
	fileCompleteTimeout = time.Second * 15
	fileCompleteSleep   = time.Millisecond * 500
)

func NewStartCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "start",
		Short: "Start erebus state streaming file proxy",
		Long: `Start erebus state streaming file proxy. The process will watch
the configured state streaming directory for changes. New files will be parsed
and sent to the configured consumer. Existing files may have been ignored, so erebus
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

	// Only notify when files are create. If we watch for Write events, then we'll
	// get an event every time a streamed file is written to. This means files
	// being currently written to when erebus starts, will be missed.
	w.FilterOps(watcher.Create)

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
		errCh <- w.Start(watcherDuration)
	}()

	logger.Info().Str("dir", stateStreamDir).Str("file_prefix", ssFilePrefix).Msg("watching state streaming directory")

	for {
		select {
		case event := <-w.Event:
			logger.Debug().Str("event", event.String()).Msg("received new watcher event")

			// Ensure event corresponds to a newly created file. Note, the file is not
			// guaranteed to be completely written to when the event is triggered, so
			// we must ensure the file is complete prior to decoding and sending to
			// the consumer.
			if event.Op == watcher.Create && !event.IsDir() {
				if err := waitForCompleteFile(event.Path); err != nil {
					logger.Error().Err(err).Str("file", event.Path).Msg("failed to wait for file to complete; skipping...")
				} else {
					switch {
					case isDataFile(event.Path):
						pairs, err := codec.ParseDataFile(event.Path)
						if err != nil {
							logger.Error().Err(err).Str("file", event.Path).Msg("failed to parse data file; skipping...")
							continue
						}

						// TODO: send to consumer
					case isMetaFile(event.Path):
						meta, err := codec.ParseMetaFile(event.Path)
						if err != nil {
							logger.Error().Err(err).Str("file", event.Path).Msg("failed to parse meta file; skipping...")
							continue
						}

						// TODO: send to consumer
					default:
						logger.Error().Str("file", event.Path).Msg("unexpected file")
					}
				}
			}

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

func isDataFile(filePath string) bool {
	return strings.HasSuffix(filePath, "-data")
}

func isMetaFile(filePath string) bool {
	return strings.HasSuffix(filePath, "-meta")
}

// waitForCompleteFile should be called when a new file event is triggered to
// ensure the file is complete before we attempt to read and parse it. An error
// is returned if the file is not considered complete within a given time
// duration.
func waitForCompleteFile(filePath string) error {
	ctx, cancel := context.WithTimeout(context.Background(), fileCompleteTimeout)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timed out waiting for file '%s' to complete", filePath)

		default:
			ok, err := io.IsFileComplete(filePath)
			if ok && err == nil {
				cancel()
				return nil
			}

			time.Sleep(fileCompleteSleep)
		}
	}
}
