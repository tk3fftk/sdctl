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
		Long:    "validate yamls, handle banners, start build from CLI",
		Version: "0.2.0",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	cmd.AddCommand(
		NewCmdBanner(api),
		NewCmdBuild(api),
		NewCmdClear(config),
		NewCmdContext(config, api),
		NewCmdGet(config, api),
		NewCmdSet(config, api),
		NewCmdValidate(api),
		NewCmdValidateTemplate(api))
	return cmd
}
