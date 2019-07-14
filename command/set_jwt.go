package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/pkg/sdctl_context"
	"github.com/tk3fftk/sdctl/util"
)

type SetJWTOption struct {
	Config sdctl_context.SdctlConfig
	API    sdapi.SDAPI
}

func NewCmdSetJWT(config sdctl_context.SdctlConfig, api sdapi.SDAPI) *cobra.Command {
	o := &SetJWTOption{
		Config: config,
		API:    api,
	}
	cmd := &cobra.Command{
		Use:   "JWT",
		Short: "set your Screwdriver.cd JWT url",
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}
	return cmd
}

func (o *SetJWTOption) Run(cmd *cobra.Command, args []string) error {
	token, err := o.API.GetJWT()
	if err != nil {
		return err
	}

	configPATH, err := util.ConfigPATH()
	if err != nil {
		return err
	}
	o.Config.SetParam(sdctl_context.SDJWTKey, token, nil)
	o.Config.Update(configPATH)
	println("Bearer " + token)
	return nil
}
