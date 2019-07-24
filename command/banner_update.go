package command

import (
	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

type BannerUpdateOption struct {
	API sdapi.SDAPI
}

var (
	id         string
	msg        string
	bannerType string
	isActive   string
	delete     bool
)

func NewCmdBannerUpdate(api sdapi.SDAPI) *cobra.Command {
	o := &BannerUpdateOption{
		API: api,
	}

	cmd := &cobra.Command{
		Use:   "set",
		Short: "update a banner",
		Long:  "update a banner with POST, PUT, DELETE method. specify banner ID when using PUT or DELETE method",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
	}

	cmd.Flags().StringVarP(&id, "id", "i", "", "specify banner ID when update or delete")
	cmd.Flags().StringVarP(&msg, "msg", "m", "", "banner message body")
	cmd.Flags().StringVarP(&bannerType, "type", "t", "info", "banner type (info, warn)")
	cmd.Flags().StringVarP(&isActive, "active", "a", "true", "banner status flag (true, false)")
	cmd.Flags().BoolVarP(&delete, "delete", "d", false, "flag for delete banner (required with id)")

	return cmd
}

func (o *BannerUpdateOption) Run(cmd *cobra.Command, args []string) error {
	if msg == "" && id == "" {
		return cmd.Help()
	}

	_, err := o.API.UpdateBanner(id, msg, bannerType, isActive, delete, false)
	if err != nil {
		return err
	}

	return nil
}
