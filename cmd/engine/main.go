package engine

import (
	"fmt"
	"os"

	config "algotrade_service/configs"
	"algotrade_service/internal/engine/controller"

	"github.com/spf13/cobra"
)

const (
	defaultTokenEnv = "TINKOFF_TOKEN"
)

type EngineCmd struct {
	cmd        *cobra.Command
	tokenEnv   string
	configPath string
}

func New() *EngineCmd {
	engineCmd := &EngineCmd{}
	engineCmd.cmd = &cobra.Command{
		Use:     "engine",
		Short:   "run engine",
		Version: "0.1",
		RunE:    engineCmd.runE,
	}
	engineCmd.cmd.Flags().StringVar(&engineCmd.tokenEnv, "token_env", defaultTokenEnv, "env where token is stored")
	engineCmd.cmd.Flags().StringVar(&engineCmd.configPath, "config", "", "config for application")

	return engineCmd
}

func (ec *EngineCmd) GetCmd() *cobra.Command {
	return ec.cmd
}

func (ec *EngineCmd) runE(cmd *cobra.Command, args []string) error {
	token := os.Getenv(ec.tokenEnv)
	if token == "" {
		return fmt.Errorf("token is empty, check your environment")
	}

	cfg, err := config.GetConfig(ec.configPath)
	if err != nil {
		return fmt.Errorf("cannot get config, err: %v", err)
	}
	// new data
	// new task
	// new controller(data, task)

	controller.NewController(token, cfg)

	return nil
}
