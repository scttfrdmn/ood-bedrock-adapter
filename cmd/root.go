// Package cmd implements the ood-bedrock-adapter CLI.
package cmd

import (
	"github.com/spf13/cobra"
)

var version = "dev" // overridden at release time via -ldflags -X .../cmd.version

var (
	region  string
	modelID string
	roleArn string
)

var rootCmd = &cobra.Command{
	Version: version,
	Use:     "ood-bedrock-adapter",
	Short:   "OOD compute adapter for AWS Bedrock batch inference",
	Long: `ood-bedrock-adapter submits, monitors, and cancels AWS Bedrock batch
model-invocation jobs (CreateModelInvocationJob) on behalf of Open OnDemand.

Batch inference fans a foundation model out across an S3 JSONL input manifest and
writes results to an S3 output prefix — a natural fit for OOD's submit/poll model
(unlike real-time single-request inference). Job specs are read from stdin as JSON.`,
}

// Execute runs the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&region, "region", "us-east-1", "AWS region")
	pf.StringVar(&modelID, "model-id", "", "Bedrock model ID (e.g. anthropic.claude-3-5-sonnet-20241022-v2:0); overridden by job spec model_id")
	pf.StringVar(&roleArn, "role-arn", "", "IAM role ARN Bedrock assumes to read/write the S3 manifests; overridden by job spec role_arn")
}
