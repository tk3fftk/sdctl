package main

import (
	"fmt"
	"os"

	"github.com/tk3fftk/sdctl/command"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"github.com/tk3fftk/sdctl/util"
)

func failureExit(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %v\n", err)
	}
	os.Exit(1)
}

func main() {

	configPATH, err := util.ConfigPATH()
	if err != nil {
		failureExit(err)
	}
	config, err := sdctl_context.LoadConfig(configPATH, false)
	if err != nil {
		failureExit(err)
	}
	sdctx := config.SdctlContexts[config.CurrentContext]
	api, err := sdapi.New(sdctx, nil)
	if err != nil {
		failureExit(err)
	}

	cmd := command.NewCmd(config, api)
	if err := cmd.Execute(); err != nil {
		failureExit(err)
	}
}
