package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

func NewCmd(config sdctl_context.SdctlConfig, api sdapi.SDAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "sdctl",
		Short:   "Screwdriver.cd API wrapper",
		Long:    "validate yamls, start build locally",
		Version: "0.1.0",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewCmdGet(config, api),
		NewCmdSet(config, api),
		NewCmdContext(config, api),
		NewCmdClear(config),
		NewCmdBuild(api))
	return cmd
}
