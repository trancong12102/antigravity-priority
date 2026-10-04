package runtime

import "testing"

func TestIsRateLimitFailure(t *testing.T) {
	cases := []struct {
		name   string
		status int
		body   string
		want   bool
	}{
		{"http 429", 429, "", true},
		{"http 429 with text body", 429, "RESOURCE_EXHAUSTED: quota rate limit exceeded", true},
		{"google error status", 0, `{"error":{"code":429,"message":"quota","status":"RESOURCE_EXHAUSTED"}}`, true},
		{"google error code only", 0, `{"error":{"code":429,"message":"slow down"}}`, true},
		{"google streaming array", 0, `[{"error":{"code":429,"status":"RESOURCE_EXHAUSTED"}}]`, true},
		{"error info reason", 0, `{"error":{"code":8,"details":[{"@type":"type.googleapis.com/google.rpc.ErrorInfo","reason":"RATE_LIMIT_EXCEEDED"}]}}`, true},
		{"anthropic type", 0, `{"type":"error","error":{"type":"rate_limit_error","message":"too many"}}`, true},
		{"openai code", 0, `{"error":{"code":"rate_limit_exceeded","message":"too many"}}`, true},
		{"json after prefix", 0, `upstream error: {"error":{"status":"RESOURCE_EXHAUSTED"}}`, true},
		{"429 inside message text", 500, `{"error":{"code":500,"message":"backend 429 retries exhausted","status":"INTERNAL"}}`, false},
		{"429 in request id", 400, `{"error":{"code":400,"message":"bad request id req-429-abc","status":"INVALID_ARGUMENT"}}`, false},
		{"plain text mentioning 429", 0, "failed after reading 4290 bytes", false},
		{"plain text mentioning rate limit", 0, "RESOURCE_EXHAUSTED rate_limit in prose", false},
		{"other http status", 503, `{"error":{"code":503,"status":"UNAVAILABLE"}}`, false},
		{"empty", 0, "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := isRateLimitFailure(tc.status, tc.body); got != tc.want {
				t.Fatalf("isRateLimitFailure(%d, %q) = %v, want %v", tc.status, tc.body, got, tc.want)
			}
		})
	}
}
