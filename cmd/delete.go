package cmd

import (
	"context"
	"fmt"

	"github.com/scttfrdmn/ood-bedrock-adapter/internal/bedrock"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete <job-arn>",
	Short: "Stop (cancel) a Bedrock batch job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := bedrock.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}
		if err := client.StopJob(ctx, args[0]); err != nil {
			return err
		}
		fmt.Printf("Job %s stopped\n", args[0])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
