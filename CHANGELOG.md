# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.1.0] - 2026-05-30

### Added
- Initial release — OOD compute adapter for AWS Bedrock batch inference, translating Open
  OnDemand job submissions to Bedrock control-plane `CreateModelInvocationJob` API calls
  (aws-openondemand#11).
- CLI commands: `submit` (JSON job spec from stdin → `CreateModelInvocationJob` → prints the
  job ARN), `status <job-arn>` (`GetModelInvocationJob` → OOD-normalized status), `delete
  <job-arn>` (`StopModelInvocationJob`), and `info <job-arn>` (full job detail as JSON).
- Root flags `--region`, `--model-id`, `--role-arn`; job-spec `model_id`/`role_arn` override
  the flags.
- Table-driven unit test for the Bedrock `ModelInvocationJobStatus` → OOD status mapping
  (all 10 SDK states + an unknown sentinel).
- Skeleton substrate integration test (`-tags=integration`), skipped pending a substrate
  Bedrock control-plane batch-inference plugin (the existing `bedrock_runtime` plugin covers
  only the data-plane `InvokeModel` API).
- CI workflow with pinned action SHAs; tag-triggered goreleaser release (cosign-signed
  multi-arch binaries).
