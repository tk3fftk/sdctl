package command

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
	"github.com/tk3fftk/sdctl/pkg/sdapi"
)

type SecretSetOption struct {
	API         sdapi.SDAPI
	PipelineID  string
	SecretKey   string
	SecretValue string
	AllowInPR   bool
}

func NewCmdSecretSet(api sdapi.SDAPI) *cobra.Command {
	o := &SecretSetOption{
		API: api,
	}
	cmd := &cobra.Command{
		Use:   "set",
		Short: "set secret to pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run(cmd, args)
		},
	}

	cmd.Flags().StringVarP(&o.PipelineID, "pipeline", "p", "", "specify pipeline id")
	_ = cmd.MarkFlagRequired("pipeline")
	cmd.Flags().StringVarP(&o.SecretKey, "key", "k", "", "SECRET_KEY")
	_ = cmd.MarkFlagRequired("key")
	cmd.Flags().StringVarP(&o.SecretValue, "value", "v", "", "SECRET_VALUE")
	_ = cmd.MarkFlagRequired("value")
	cmd.Flags().BoolVarP(&o.AllowInPR, "allow-in-pr", "", false, "ALLOW_IN_PR")

	return cmd
}

func (o *SecretSetOption) Run(cmd *cobra.Command, args []string) error {
	pipelineIDNum, err := strconv.Atoi(o.PipelineID)
	if err != nil {
		return fmt.Errorf("failed to convert %s to int: %v", o.PipelineID, err)
	}

	// Screwdriver allow only "/^[A-Z_][A-Z0-9_]*$/]" as secret key
	uppperKey := strings.ToUpper(o.SecretKey)
	if err := o.API.UpdateSecret(pipelineIDNum, uppperKey, o.SecretValue, o.AllowInPR); err != nil {
		return fmt.Errorf("failed to set secret: %v", err)
	}

	fmt.Fprintf(os.Stdout, "setting secret %s is succeed!\n", uppperKey)
	return nil
}
