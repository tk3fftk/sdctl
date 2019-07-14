package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

type BuildOption struct {
	API sdapi.SDAPI
}

func NewCmdBuild(api sdapi.SDAPI) *cobra.Command {
	o := &BuildOption{
		API: api,
	}
	cmd := &cobra.Command{
		Use:     "build <pipelieid> <start_from>",
		Short:   "start a job.",
		Aliases: []string{"b"},
		Run: func(cmd *cobra.Command, args []string) {
			o.Run(cmd, args)
		},
	}
	return cmd
}

func (o *BuildOption) Run(cmd *cobra.Command, args []string) error {
	if len(args) != 2 {
		return cmd.Help()
	}
	pipelineID := args[0]
	startFrom := args[1]
	if err := o.API.PostEvent(pipelineID, startFrom, false); err != nil {
		return err
	}
	return nil
}
