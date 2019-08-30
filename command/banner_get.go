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
	o.printColumn("ID", "IsActive", "Message")
	for _, b := range banners {
		o.printColumn(b.ID, b.IsActive, b.Message)
	}
}

func (o *BannerGetOption) printColumn(id, isActive, msg interface{}) {
	fmt.Fprintf(os.Stdout, "%-8v%-12v%v\n", id, isActive, msg)
}
