package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"github.com/tk3fftk/sdctl/util"
)

type SetAPIOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdSetAPI(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &SetAPIOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "api",
		Short: "set your Screwdriver.cd api url",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage: true,
	}
	return cmd
}

func (o *SetAPIOption) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}

	configPATH, err := util.ConfigPATH()
	if err != nil {
		return err
	}
	o.Config.SetParam(sdctl_context.APIURLKey, args[0], nil)
	o.Config.Update(configPATH)
	return nil
}
