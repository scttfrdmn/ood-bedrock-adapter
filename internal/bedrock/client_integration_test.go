//go:build integration

package bedrock_test

import (
	"testing"
)

// TestCreateGetStopJob_Substrate will exercise the full Bedrock batch job
// lifecycle (CreateModelInvocationJob → GetModelInvocationJob →
// StopModelInvocationJob) against the substrate emulator.
//
// It is skipped until substrate ships a Bedrock control-plane batch-inference
// plugin. The existing substrate bedrock_runtime plugin covers only the data
// plane (InvokeModel/ApplyGuardrail), not the control-plane ModelInvocationJob
// APIs this adapter uses. Tracking: scttfrdmn/substrate#297.
func TestCreateGetStopJob_Substrate(t *testing.T) {
	t.Skip("pending substrate Bedrock control-plane batch-inference plugin; see scttfrdmn/substrate#297")
}
