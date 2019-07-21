package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
	"github.com/tk3fftk/sdctl/util"
)

type ValidateOption struct {
	API sdapi.SDAPI
}

var pipelineFilePATH string
var validatedOutput bool

func NewCmdValidate(api sdapi.SDAPI) *cobra.Command {
	o := &ValidateOption{
		API: api,
	}
	cmd := &cobra.Command{
		Use:     "validate",
		Short:   "validate your screwdriver.yaml, default to screwdriver.yaml",
		Aliases: []string{"v"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().StringVarP(&pipelineFilePATH, "file", "f", "screwdriver.yaml", "specify pipeline file path")
	cmd.Flags().BoolVarP(&validatedOutput, "output", "o", false, "print velidator result")

	return cmd
}

func (o *ValidateOption) Run(cmd *cobra.Command, args []string) error {
	yaml, err := util.ReadYaml(pipelineFilePATH)
	if err != nil {
		return err
	}
	if err := o.API.Validator(yaml, false, validatedOutput); err != nil {
		return err
	}
	return nil
}
