package cmd

import (
	"context"
	"encoding/json"
	"os"

	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrock/types"
	"github.com/scttfrdmn/ood-bedrock-adapter/internal/bedrock"
	internalood "github.com/scttfrdmn/ood-bedrock-adapter/internal/ood"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status <job-arn>",
	Short: "Get the status of a Bedrock batch job",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		client, err := bedrock.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}

		detail, err := client.GetJob(ctx, args[0])
		if err != nil {
			return err
		}

		js := internalood.JobStatus{
			ID:     args[0],
			Status: bedrockStateToOod(detail.Status),
		}
		if detail.Message != nil {
			js.Message = *detail.Message
		}

		return json.NewEncoder(os.Stdout).Encode(js)
	},
}

// bedrockStateToOod maps a Bedrock ModelInvocationJobStatus to an OOD status string.
func bedrockStateToOod(s bedrocktypes.ModelInvocationJobStatus) string {
	switch s {
	case bedrocktypes.ModelInvocationJobStatusSubmitted,
		bedrocktypes.ModelInvocationJobStatusValidating,
		bedrocktypes.ModelInvocationJobStatusScheduled:
		return internalood.StatusQueued
	case bedrocktypes.ModelInvocationJobStatusInProgress,
		bedrocktypes.ModelInvocationJobStatusStopping:
		return internalood.StatusRunning
	case bedrocktypes.ModelInvocationJobStatusCompleted,
		bedrocktypes.ModelInvocationJobStatusPartiallyCompleted:
		return internalood.StatusCompleted
	case bedrocktypes.ModelInvocationJobStatusFailed,
		bedrocktypes.ModelInvocationJobStatusExpired:
		return internalood.StatusFailed
	case bedrocktypes.ModelInvocationJobStatusStopped:
		return internalood.StatusCancelled
	default:
		return internalood.StatusUnknown
	}
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
