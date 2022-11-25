package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	blocks  uint64
	appName = "cosi"
	unsafe  bool
)

// NewRootCmd returns the root command for relayer.

func NewRootCmd() *cobra.Command {
	// RootCmd represents the base command when called without any subcommands
	var rootCmd = &cobra.Command{
		Use:   appName,
		Short: "cosi is a tool to reindex tendermint events if they are missed",
	}

	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, _ []string) error {
		// reads `homeDir/config.yaml` into `var config *Config` before each command
		// if err := initConfig(rootCmd); err != nil {
		// 	return err
		// }

		return nil
	}

	// --unsafe flag
	rootCmd.PersistentFlags().BoolVar(&unsafe, "unsafe", true, "dont wait on delivery confirmation")
	if err := viper.BindPFlag("unsafe", rootCmd.PersistentFlags().Lookup("unsafe")); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(
		rpcCmd(),
		kvCmd(),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.

func Execute() {
	cobra.EnableCommandSorting = false

	rootCmd := NewRootCmd()
	rootCmd.SilenceUsage = true
	rootCmd.CompletionOptions.DisableDefaultCmd = true

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
