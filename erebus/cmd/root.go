package cmd

import (
	"github.com/spf13/cobra"
)

const (
	flagConfig            = "config"
	flagStateStreamingDir = "state-streaming-dir"
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

	return cmd
}
