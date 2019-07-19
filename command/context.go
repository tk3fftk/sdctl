package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

func NewCmdContext(config sdctl_context.SdctlConfig, api sdapi.SDAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "handle screwdriver contexts",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(
		NewCmdContextList(config),
		NewCmdContextCurrent(config),
		NewCmdContextSet(config))
	return cmd
}
