package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

type GetJWTOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdGetJWT(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &GetJWTOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "jwt",
		Short: "show your jwt",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage: true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *GetJWTOption) Run(cmd *cobra.Command, args []string) error {
	o.Config.PrintParam(sdctl_context.SDJWTKey, nil)
	return nil
}
