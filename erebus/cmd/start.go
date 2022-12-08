package cmd

import (
	"fmt"
	"os"

	"github.com/numiadata/tools/erebus/config"
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
	cfgFile, err := cmd.Flags().GetString(flagConfig)
	if err != nil {
		return err
	}

	cfg, err := config.Parse(cfgFile)
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	stateStreamDir, err := cmd.Flags().GetString(flagStateStreamingDir)
	if err != nil {
		return err
	}

	if _, err := os.Stat(stateStreamDir); os.IsNotExist(err) {
		return fmt.Errorf("state streaming directory '%s' does not exist", stateStreamDir)
	}

	return nil
}
