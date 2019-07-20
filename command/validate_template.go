package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/util"
)

type ValidateTemplateOption struct {
	API sdapi.SDAPI
}

var templateFilePATH string

func NewCmdValidateTemplate(api sdapi.SDAPI) *cobra.Command {
	o := &ValidateTemplateOption{
		API: api,
	}
	cmd := &cobra.Command{
		Use:     "validate-template",
		Short:   "validate your sd-template.yaml, default to sd-template.yaml",
		Aliases: []string{"vt"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage: true,
		SilenceErrors: true,
	}
	cmd.Flags().StringVarP(&templateFilePATH, "file", "f", "sd-template.yaml", "specify template file path")
	return cmd
}

func (o *ValidateTemplateOption) Run(cmd *cobra.Command, args []string) error {
	yaml, err := util.ReadYaml(templateFilePATH)
	if err != nil {
		return err
	}
	if err := o.API.ValidatorTemplate(yaml, false); err != nil {
		return err
	}
	return nil
}
