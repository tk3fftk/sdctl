package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"github.com/tk3fftk/sdctl/util"
)

type ClearOption struct {
	Config sdctl_context.SdctlConfig
}

func NewCmdClear(config sdctl_context.SdctlConfig) *cobra.Command {
	o := &ClearOption{
		Config: config,
	}
	cmd := &cobra.Command{
		Use:   "clear",
		Short: "clear your setting and set to default",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *ClearOption) Run(cmd *cobra.Command, args []string) error {
	configPATH, err := util.ConfigPATH()
	if err != nil {
		return err
	}
	sdctl_context.LoadConfig(configPATH, true)
	return nil
}
