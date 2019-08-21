package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

func NewCmdSecret(api sdapi.SDAPI) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "secret",
		Short: "handle screwdriver secrets (write only)",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
		Aliases: []string{"sec"},
	}

	cmd.AddCommand(
		NewCmdSecretSet(api),
	)
	return cmd
}
