package cli

import (
	"github.com/ondrovic/bambulab-authenticator/internal/types"
	"github.com/spf13/cobra"
)

var (
	Options = types.CliFlags{}
	RootCmd = &cobra.Command{
		Use:   "bambulab-authenticator",
		Short: "A CLI tool to export authentication info to a json file",
	}
)

func InitializeCommands() {
	initAuthenticateFlags()
	RootCmd.AddCommand(authenticateCmd)
}

func Execute() error {
	if err := RootCmd.Execute(); err != nil {
		return err
	}

	return nil
}
