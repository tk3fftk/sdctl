package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"github.com/tk3fftk/sdctl/util"
)

type SetTokenOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdSetToken(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &SetTokenOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "token",
		Short: "set your user token",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage: true,
	}
	return cmd
}

func (o *SetTokenOption) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	configPATH, err := util.ConfigPATH()
	if err != nil {
		return err
	}
	o.Config.SetParam(sdctl_context.UserTokenKey, args[0], nil)
	o.Config.Update(configPATH)
	return nil
}
