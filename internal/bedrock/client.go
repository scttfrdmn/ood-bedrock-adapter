// Package bedrock wraps the AWS Bedrock control-plane batch-inference API
// (CreateModelInvocationJob / GetModelInvocationJob / StopModelInvocationJob)
// for the OOD adapter. This is the control plane (service/bedrock), not the
// data-plane bedrockruntime InvokeModel API.
package bedrock

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrock"
	"github.com/aws/aws-sdk-go-v2/service/bedrock/types"
)

// Client wraps the AWS Bedrock control-plane client.
type Client struct {
	svc    *bedrock.Client
	region string
}

// New creates a Bedrock client using the default AWS credential chain.
func New(ctx context.Context, region string, optFns ...func(*config.LoadOptions) error) (*Client, error) {
	opts := append([]func(*config.LoadOptions) error{config.WithRegion(region)}, optFns...)
	cfg, err := config.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return nil, fmt.Errorf("load AWS config: %w", err)
	}
	return &Client{svc: bedrock.NewFromConfig(cfg), region: region}, nil
}

// JobSpec holds the parameters for a Bedrock batch model-invocation job.
type JobSpec struct {
	ModelID  string
	InputS3  string // s3://bucket/input-manifest.jsonl
	OutputS3 string // s3://bucket/output-prefix/
	RoleArn  string
	JobName  string
}

// CreateJob submits a Bedrock batch model-invocation job and returns its ARN.
func (c *Client) CreateJob(ctx context.Context, spec JobSpec) (string, error) {
	input := &bedrock.CreateModelInvocationJobInput{
		JobName: aws.String(spec.JobName),
		ModelId: aws.String(spec.ModelID),
		RoleArn: aws.String(spec.RoleArn),
		InputDataConfig: &types.ModelInvocationJobInputDataConfigMemberS3InputDataConfig{
			Value: types.ModelInvocationJobS3InputDataConfig{
				S3Uri: aws.String(spec.InputS3),
			},
		},
		OutputDataConfig: &types.ModelInvocationJobOutputDataConfigMemberS3OutputDataConfig{
			Value: types.ModelInvocationJobS3OutputDataConfig{
				S3Uri: aws.String(spec.OutputS3),
			},
		},
	}

	out, err := c.svc.CreateModelInvocationJob(ctx, input)
	if err != nil {
		return "", fmt.Errorf("bedrock CreateModelInvocationJob: %w", err)
	}
	return aws.ToString(out.JobArn), nil
}

// GetJob returns the current detail of a Bedrock batch job (by ARN or ID).
func (c *Client) GetJob(ctx context.Context, jobID string) (*bedrock.GetModelInvocationJobOutput, error) {
	out, err := c.svc.GetModelInvocationJob(ctx, &bedrock.GetModelInvocationJobInput{
		JobIdentifier: aws.String(jobID),
	})
	if err != nil {
		return nil, fmt.Errorf("bedrock GetModelInvocationJob: %w", err)
	}
	return out, nil
}

// StopJob cancels an in-flight Bedrock batch job.
func (c *Client) StopJob(ctx context.Context, jobID string) error {
	_, err := c.svc.StopModelInvocationJob(ctx, &bedrock.StopModelInvocationJobInput{
		JobIdentifier: aws.String(jobID),
	})
	if err != nil {
		return fmt.Errorf("bedrock StopModelInvocationJob: %w", err)
	}
	return nil
}
