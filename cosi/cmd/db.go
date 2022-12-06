package cmd

import (
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/numiadata/tools/pubsub"

	"github.com/numiadata/tools/cosi/utils/kv"
	"github.com/numiadata/tools/cosi/utils/state"
)

// this command is used for reinstalling the events using a local db
// load db
// load app store and prune
// if immutable tree is not deletable we should import and export current state
// add flags for block events only
func kvCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kv [start_height] [end_height] [path_to_db]",
		Short: "reindex via the db from a start height to an end height, note this only works for txs currently",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()

			start, err := strconv.ParseInt(args[0], 10, 0)
			if err != nil {
				return err
			}

			end, err := strconv.ParseInt(args[1], 10, 0)
			if err != nil {
				return err
			}

			consumer, err := pubsub.NewEventSink()
			if err != nil {
				return err
			}

			// loop through specified heights and index
			return kv.IndexTxs(ctx, consumer, args[2], start, end, unsafe)
		},
	}
	return cmd
}

// this command is used for reinstalling the events using a local db
// load db
// load app store and prune
// if immutable tree is not deletable we should import and export current state
// add flags for block events only
func stateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "state [start_height] [end_height] [path_to_db]",
		Short: "reindex via the state db from a start height to an end height, note this only works for txs currently",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			ctx := cmd.Context()

			start, err := strconv.ParseInt(args[0], 10, 0)
			if err != nil {
				return err
			}

			end, err := strconv.ParseInt(args[1], 10, 0)
			if err != nil {
				return err
			}

			consumer, err := pubsub.NewEventSink()
			if err != nil {
				return err
			}

			// loop through specified heights and index
			return state.Index(ctx, consumer, args[2], start, end, unsafe)
		},
	}
	return cmd
}

func baseHeightCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "base [path_to_db]",
		Short: "Get the base and highest height of the db",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			base, height, err := state.GetBaseHeight(args[0])
			if err != nil {
				return err
			}

			fmt.Println("base: ", base)
			fmt.Println("height: ", height)

			return nil
		},
	}
	return cmd
}

func compactStateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "compact [path_to_db]",
		Short: "compact the diskspace of the state db",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			return state.ForceCompact(args[0])
		},
	}
	return cmd
}
