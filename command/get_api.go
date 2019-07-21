package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

type GetAPIOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdGetAPI(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &GetAPIOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "api",
		Short: "get configured api url",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *GetAPIOption) Run(cmd *cobra.Command, args []string) error {
	o.Config.PrintParam(sdctl_context.APIURLKey, nil)
	return nil
}
