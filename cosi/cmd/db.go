package cmd

import (
	"strconv"

	"github.com/spf13/cobra"

	"github.com/numiadata/tools/cosi/utils/kv"
	"github.com/numiadata/tools/cosi/utils/pubsub"
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
