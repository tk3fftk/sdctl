package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

func NewCmdGet(config sdctl_context.SdctlConfig, api sdapi.SDAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get sdctl settings and Screwdriver.cd information",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewCmdGetToken(config),
		NewCmdGetAPI(config),
		NewCmdGetJWT(config),
		NewCmdGetBuildPages(api))
	return cmd
}
