package cli

import (
	"fmt"

	"github.com/ondrovic/bambulab-authenticator/internal/auth"
	"github.com/ondrovic/bambulab-authenticator/internal/consts"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	authenticateCmd = &cobra.Command{
		Use:   "authenticate",
		Short: "Authenticate with your credentials",
		Args:  cobra.ExactArgs(0),
		RunE:  runAuthenticate,
	}
)

func initAuthenticateFlags() {

	authenticateCmd.Flags().StringVarP(&Options.OutputPath, "output-path", "o", consts.EMPTY_STRING, "Output path of the authentication info")
	authenticateCmd.Flags().StringVarP(&Options.UserAccount, "user-account", "u", consts.EMPTY_STRING, "User account")
	authenticateCmd.Flags().StringVarP(&Options.UserPassword, "user-password", "p", consts.EMPTY_STRING, "User account password")
	authenticateCmd.Flags().StringVarP(&Options.UserRegion, "user-region", "r", consts.EMPTY_STRING, "User region")

	markAllFlagsRequired(authenticateCmd)

	viper.BindPFlags(authenticateCmd.Flags())
}

func markAllFlagsRequired(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(flag *pflag.Flag) {
		err := cmd.MarkFlagRequired(flag.Name)
		if err != nil {
			fmt.Printf("error setting flag: %s required: %v", flag.Name, err)
		}
	})
}

func runAuthenticate(cmd *cobra.Command, args []string) error {

	if err := auth.Login(&Options); err != nil {
		return err
	}

	return nil
}
