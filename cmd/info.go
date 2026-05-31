package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/scttfrdmn/ood-bedrock-adapter/internal/bedrock"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info <job-arn>",
	Short: "Print full Bedrock batch job details as JSON",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := bedrock.New(ctx, region)
		if err != nil {
			return err
		}
		detail, err := client.GetJob(ctx, args[0])
		if err != nil {
			return err
		}
		return json.NewEncoder(os.Stdout).Encode(detail)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}
