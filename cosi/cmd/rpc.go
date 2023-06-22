package cmd

import (
	"context"

	"github.com/cometbft/cometbft/rpc/client/http"
	"github.com/spf13/cobra"
)

// load db
// load app store and prune
// if immutable tree is not deletable we should import and export current state
// add flags for block events only
func rpcCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "rpc [start height] [end height] [rpc address]",
		Short: "reindex via the rpc from a start height to an end height",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {

			// start, err := strconv.ParseInt(args[0], 10, 0)
			// if err != nil {
			// 	return err
			// }

			// end, err := strconv.ParseInt(args[1], 10, 0)
			// if err != nil {
			// 	return err
			// }

			// reindexBlock(ctx, start, end, args[2])
			return nil
		},
	}
	return cmd
}

func reindexBlock(ctx context.Context, start, end int64, rpcAddress string) error {
	// setup pubsub

	// create loop from start height to end height
	//    should we run concurrently?
	for i := start; i < end; i++ {

		c, err := http.New(rpcAddress, "/websocket")
		if err != nil {
			return err
		}

		// query for block_results at height
		c.BlockResults(context.Background(), &i)
		// if err != nil {
		// 	// print error for easier debugging
		// 	fmt.Println(err)
		// 	return err
		// }
		// get all the events from the block

		// send to pubsub
	}
	return nil
}
