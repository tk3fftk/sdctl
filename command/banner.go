package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

func NewCmdBanner(api sdapi.SDAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "banner",
		Short:   "handle screwdriver banners",
		Aliases: []string{"bn"},
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}

	cmd.AddCommand(NewCmdBannerGet(api))
	cmd.AddCommand(NewCmdBannerUpdate(api))

	return cmd
}
