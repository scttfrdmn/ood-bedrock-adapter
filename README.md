# ood-bedrock-adapter

[![CI](https://github.com/scttfrdmn/ood-bedrock-adapter/actions/workflows/ci.yml/badge.svg)](https://github.com/scttfrdmn/ood-bedrock-adapter/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/scttfrdmn/ood-bedrock-adapter)](https://goreportcard.com/report/github.com/scttfrdmn/ood-bedrock-adapter)
[![codecov](https://codecov.io/gh/scttfrdmn/ood-bedrock-adapter/branch/main/graph/badge.svg)](https://codecov.io/gh/scttfrdmn/ood-bedrock-adapter)
[![Go Reference](https://pkg.go.dev/badge/github.com/scttfrdmn/ood-bedrock-adapter.svg)](https://pkg.go.dev/github.com/scttfrdmn/ood-bedrock-adapter)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

OOD compute adapter for **AWS Bedrock batch inference**. Translates Open OnDemand job
submissions to AWS Bedrock control-plane `CreateModelInvocationJob` API calls.

Real-time Bedrock inference (a single request/response) doesn't fit Open OnDemand's
submit/poll model — but **batch inference** does: point it at an S3 JSONL manifest, Bedrock
fans the foundation model out across every record and writes results to an S3 output prefix.
Useful for large-scale model-evaluation benchmarks, LLM-based annotation/literature mining,
and embedding/RAG dataset preparation.

This is the **control plane** (`service/bedrock` — `*ModelInvocationJob`), not the data-plane
`bedrockruntime` `InvokeModel` API.

## Commands

| Subcommand | Bedrock API | Description |
| --- | --- | --- |
| `submit` | `CreateModelInvocationJob` | reads a JSON job spec on stdin, prints the job ARN |
| `status <job-arn>` | `GetModelInvocationJob` | OOD-normalized status JSON |
| `delete <job-arn>` | `StopModelInvocationJob` | cancel an in-flight job |
| `info <job-arn>` | `GetModelInvocationJob` | full job detail as JSON |

## Flags

`--region` (default `us-east-1`), `--model-id` (e.g. `anthropic.claude-3-5-sonnet-20241022-v2:0`),
`--role-arn` (the role Bedrock assumes to read/write the S3 manifests). The job spec's
`model_id` / `role_arn` override the flags.

## Job spec (stdin JSON)

```json
{
  "job_name": "eval-run-2026-05",
  "model_id": "anthropic.claude-3-5-sonnet-20241022-v2:0",
  "input_s3": "s3://my-bucket/batch/input-manifest.jsonl",
  "output_s3": "s3://my-bucket/batch/output/",
  "role_arn": "arn:aws:iam::123456789012:role/bedrock-batch"
}
```

### Input manifest format (JSONL, one record per line)

```json
{"recordId": "1", "modelInput": {"prompt": "Summarize: ...", "max_tokens": 512}}
{"recordId": "2", "modelInput": {"prompt": "Classify: ...", "max_tokens": 128}}
```

## Status mapping

| Bedrock `ModelInvocationJobStatus` | OOD status |
| --- | --- |
| Submitted, Validating, Scheduled | queued |
| InProgress, Stopping | running |
| Completed, PartiallyCompleted | completed |
| Failed, Expired | failed |
| Stopped | cancelled |
| *(unknown)* | undetermined |

## Testing

Unit tests cover the status-state mapping (`go test ./...`). A substrate-backed integration
test (`-tags=integration`) is present but skipped pending a substrate Bedrock control-plane
batch-inference plugin — substrate currently emulates only the data-plane `bedrock_runtime`
(`InvokeModel`) API. See the linked substrate feature request.

## License

MIT — see [LICENSE](LICENSE).
