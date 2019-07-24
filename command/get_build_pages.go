package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

type GetBuildPagesOption struct {
	API sdapi.SDAPI
}

func NewCmdGetBuildPages(api sdapi.SDAPI) *cobra.Command {
	o := &GetBuildPagesOption{
		API: api,
	}
	cmd := &cobra.Command{
		Use:     "build-pages <BUILD_ID>",
		Short:   "get build page url",
		Aliases: []string{"bp"},
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
	}
	return cmd
}

func (o *GetBuildPagesOption) Run(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return cmd.Help()
	}
	buildID := args[0]

	if err := o.API.GetPipelinePageFromBuildID(buildID); err != nil {
		return err
	}
	return nil
}
