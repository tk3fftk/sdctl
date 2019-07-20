package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

type GetTokenOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdGetToken(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &GetTokenOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "token",
		Short: "get your user token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage: true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *GetTokenOption) Run(cmd *cobra.Command, args []string) error {
	o.Config.PrintParam(sdctl_context.UserTokenKey, nil)
	return nil
}
