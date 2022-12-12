package cmd_test

import (
	"context"
	"io"
	"os"
	"path"
	"testing"
	"time"

	"github.com/numiadata/tools/erebus/cmd"
	"github.com/numiadata/tools/erebus/consumer"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"

	storetypes "github.com/cosmos/cosmos-sdk/store/types"
)

func TestWatchStreamingDir(t *testing.T) {
	tmpDir := t.TempDir()
	ctx, cancel := context.WithCancel(context.TODO())

	tc := consumer.NewTestConsumer()

	// start watcher in a goroutine
	go func() {
		_ = cmd.WatchStreamingDir(ctx, zerolog.Nop(), tmpDir, cmd.DefaultFilePrefix, tc)
	}()

	// Allow ample time for watcher to start listening for events prior to copying
	// files.
	time.Sleep(time.Second)

	// Copy an example data and meta file from examples which the watcher should
	// catch and push messages for to the TestConsumer.
	copyFile(t, "../examples/state_streaming/block-271-data", path.Join(tmpDir, "block-271-data"))
	copyFile(t, "../examples/state_streaming/block-288-meta", path.Join(tmpDir, "block-288-meta"))

	doneCh := make(chan bool)
	go func() {
		ticker := time.NewTicker(time.Millisecond * 500)
		defer ticker.Stop()

		for range ticker.C {
			if len(tc.DataMessages()) == 1 && len(tc.MetaMessages()) == 1 {
				doneCh <- true
			}
		}
	}()

	// wait for all messages to be received
	<-doneCh

	// cancel context which should allow the WatchStreamingDir goroutine to exit
	cancel()
	<-ctx.Done()

	require.Len(t, tc.DataMessages(), 1)
	require.Len(t, tc.MetaMessages(), 1)
	require.Len(t, tc.DataMessages()[0].Body.([]storetypes.StoreKVPair), 12)
}

func copyFile(t *testing.T, src, dst string) {
	original, err := os.Open(src)
	require.NoError(t, err)

	defer original.Close()

	new, err := os.Create(dst)
	require.NoError(t, err)

	defer new.Close()

	_, err = io.Copy(new, original)
	require.NoError(t, err)
}
