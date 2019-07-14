package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
)

type ContextCurrentOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdContextCurrent(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &ContextCurrentOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "current",
		Short: "show current context",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}
	return cmd
}

func (o *ContextCurrentOption) Run(cmd *cobra.Command, args []string) error {
	o.Config.PrintParam(sdctl_context.CurrentContextKey, nil)
	return nil
}
