package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/scttfrdmn/ood-bedrock-adapter/internal/bedrock"
	"github.com/spf13/cobra"
)

// JobSpec is the Bedrock batch-inference job submission payload (from stdin).
type JobSpec struct {
	ModelID  string `json:"model_id,omitempty"` // overrides --model-id
	InputS3  string `json:"input_s3"`           // s3://bucket/input-manifest.jsonl
	OutputS3 string `json:"output_s3"`          // s3://bucket/output-prefix/
	RoleArn  string `json:"role_arn,omitempty"` // overrides --role-arn
	JobName  string `json:"job_name"`
}

var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "Submit an OOD job to AWS Bedrock batch inference",
	Long:  "Reads a JSON job spec from stdin and submits it as a Bedrock model-invocation job.",
	RunE: func(cmd *cobra.Command, args []string) error {
		var spec JobSpec
		if err := json.NewDecoder(os.Stdin).Decode(&spec); err != nil {
			return fmt.Errorf("decode job spec: %w", err)
		}

		if spec.InputS3 == "" {
			return fmt.Errorf("job spec must include input_s3")
		}
		if spec.OutputS3 == "" {
			return fmt.Errorf("job spec must include output_s3")
		}
		if spec.JobName == "" {
			return fmt.Errorf("job spec must include job_name")
		}

		effectiveModel := spec.ModelID
		if effectiveModel == "" {
			effectiveModel = modelID
		}
		if effectiveModel == "" {
			return fmt.Errorf("--model-id is required (or set model_id in job spec)")
		}

		effectiveRole := spec.RoleArn
		if effectiveRole == "" {
			effectiveRole = roleArn
		}
		if effectiveRole == "" {
			return fmt.Errorf("--role-arn is required (or set role_arn in job spec)")
		}

		ctx := context.Background()
		client, err := bedrock.New(ctx, region, awsOptions(ctx)...)
		if err != nil {
			return err
		}

		jobArn, err := client.CreateJob(ctx, bedrock.JobSpec{
			ModelID:  effectiveModel,
			InputS3:  spec.InputS3,
			OutputS3: spec.OutputS3,
			RoleArn:  effectiveRole,
			JobName:  spec.JobName,
		})
		if err != nil {
			return err
		}

		fmt.Println(jobArn)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)
}
