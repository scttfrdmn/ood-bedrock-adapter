package cmd

import (
	"testing"

	bedrocktypes "github.com/aws/aws-sdk-go-v2/service/bedrock/types"
)

func TestBedrockStateToOod(t *testing.T) {
	tests := []struct {
		state    bedrocktypes.ModelInvocationJobStatus
		expected string
	}{
		{bedrocktypes.ModelInvocationJobStatusSubmitted, "queued"},
		{bedrocktypes.ModelInvocationJobStatusValidating, "queued"},
		{bedrocktypes.ModelInvocationJobStatusScheduled, "queued"},
		{bedrocktypes.ModelInvocationJobStatusInProgress, "running"},
		{bedrocktypes.ModelInvocationJobStatusStopping, "running"},
		{bedrocktypes.ModelInvocationJobStatusCompleted, "completed"},
		{bedrocktypes.ModelInvocationJobStatusPartiallyCompleted, "completed"},
		{bedrocktypes.ModelInvocationJobStatusFailed, "failed"},
		{bedrocktypes.ModelInvocationJobStatusExpired, "failed"},
		{bedrocktypes.ModelInvocationJobStatusStopped, "cancelled"},
		{bedrocktypes.ModelInvocationJobStatus("UNKNOWN_STATE"), "undetermined"},
	}

	for _, tt := range tests {
		t.Run(string(tt.state), func(t *testing.T) {
			got := bedrockStateToOod(tt.state)
			if got != tt.expected {
				t.Errorf("bedrockStateToOod(%q) = %q, want %q", tt.state, got, tt.expected)
			}
		})
	}
}
