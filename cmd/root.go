// Package cmd implements the ood-bedrock-adapter CLI.
package cmd

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/scttfrdmn/ood-bedrock-adapter/internal/awscfg"
	"github.com/spf13/cobra"
)

var (
	assumeRoleArn     string // #78: per-user role to assume (empty = instance role)
	assumeRoleExtID   string
	assumeRoleSession string
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

// #78: per-user cross-account AssumeRole flags. Empty (default) = use the OOD instance role.
func init() {
	pf := rootCmd.PersistentFlags()
	pf.StringVar(&assumeRoleArn, "assume-role-arn", "", "IAM role ARN to assume for AWS calls (empty = use the instance role)")
	pf.StringVar(&assumeRoleExtID, "assume-role-external-id", "", "sts:ExternalId for the assumed-role trust policy")
	pf.StringVar(&assumeRoleSession, "assume-role-session-name", "", "RoleSessionName for the assumed role (e.g. the OOD username)")
}

// awsOptions builds the AWS config options from the root flags (region + optional AssumeRole).
func awsOptions(ctx context.Context) []func(*config.LoadOptions) error {
	return awscfg.LoadOptions(ctx, awscfg.Options{
		Region:        region,
		AssumeRoleARN: assumeRoleArn,
		ExternalID:    assumeRoleExtID,
		SessionName:   assumeRoleSession,
	})
}
