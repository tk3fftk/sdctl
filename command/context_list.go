package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

type ContextListOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdContextList(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &ContextListOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:     "list",
		Short:   "show context list",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *ContextListOption) Run(cmd *cobra.Command, args []string) error {
	o.Config.PrintParam(sdctl_context.ContextsKey, nil)
	return nil
}
