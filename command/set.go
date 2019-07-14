package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

func NewCmdSet(config sdctl_context.SdctlConfig, api sdapi.SDAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set sdctl settings",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewCmdSetToken(config),
		NewCmdSetAPI(config),
		NewCmdSetJWT(config, api))
	return cmd
}
