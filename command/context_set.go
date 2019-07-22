package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"github.com/tk3fftk/sdctl/util"
)

type ContextSetOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdContextSet(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &ContextSetOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "set <context>",
		Short: "set current to context. if it doesn't exist, create new one",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *ContextSetOption) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	context := args[0]
	configPATH, err := util.ConfigPATH()
	if err != nil {
		return err
	}
	o.Config.SetParam(sdctl_context.CurrentContextKey, context, nil)
	o.Config.Update(configPATH)
	return nil
}
