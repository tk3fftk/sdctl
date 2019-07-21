package command

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

type BannerGetOption struct {
	API sdapi.SDAPI
}

func NewCmdBannerGet(api sdapi.SDAPI) *cobra.Command {
	o := &BannerGetOption{
		API: api,
	}
	cmd := &cobra.Command{
		Use:   "get",
		Short: "get a list of banners",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	return cmd
}

func (o *BannerGetOption) Run(cmd *cobra.Command, args []string) error {
	if len(args) > 0 {
		return cmd.Help()
	}

	banners, err := o.API.GetBanners()
	if err != nil {
		return err
	}
	o.print(banners)

	return nil
}

func (o *BannerGetOption) print(banners []sdapi.BannerResponse) {
	fmt.Fprintf(os.Stdout, "ID\tIsActive\tMessage\n")
	for _, b := range banners {
		fmt.Fprintf(os.Stdout, "%v\t%v\t%v\n", b.ID, b.IsActive, b.Message)
	}
}
