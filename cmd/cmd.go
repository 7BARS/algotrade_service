package cmd

import (
	"github.com/spf13/cobra"

	"algotrade_service/cmd/engine"
)

func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:     "version",
		Short:   "run algotrade service",
		Version: "0.1",
	}
	engineCmd := engine.New()
	rootCmd.AddCommand(engineCmd.GetCmd())

	return rootCmd
}
