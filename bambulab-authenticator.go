package main

import (
	"github.com/ondrovic/bambulab-authenticator/cmd/cli"
	sCli "github.com/ondrovic/common/utils/cli"
	"runtime"
)

func main() {
	if err := sCli.ClearTerminalScreen(runtime.GOOS); err != nil {
		return
	}

	cli.InitializeCommands()

	if err := cli.RootCmd.Execute(); err != nil {
		return
	}
}
