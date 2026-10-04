package runtime_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"antigravity-priority/internal/host"
	"antigravity-priority/internal/runtime"
)

func newCooldownCandidate(t *testing.T, configYAML string) (*runtime.Runtime, string) {
	t.Helper()
	authPath := filepath.Join(t.TempDir(), "auth.json")
	if err := os.WriteFile(authPath, []byte(`{"priority":100,"disabled":false}`), 0o600); err != nil {
		t.Fatal(err)
	}
	mock := newMockHost()
	mock.files = []host.AuthFile{{Name: "account", AuthIndex: "auth", Provider: "antigravity", Priority: 100}}
	mock.authDocs["auth"] = host.AuthDocument{AuthIndex: "auth", Path: authPath}
	r := newTestRuntime(t, runtime.Options{Host: mock})
	cachePath := filepath.Join(t.TempDir(), "cache.json")
	if _, err := r.Register(context.Background(), runtime.RegisterRequest{ConfigYAML: configYAML + "state_cache_path: " + filepath.ToSlash(cachePath) + "\n"}); err != nil {
		t.Fatal(err)
	}
	return r, authPath
}

func readPriorityDocument(t *testing.T, path string) string {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestRuntime_429CooldownDoesNotWriteWhenPluginDisabled(t *testing.T) {
	r, authPath := newCooldownCandidate(t, "enabled: false\n")
	payload := []byte(`{"Provider":"antigravity","AuthIndex":"auth","Failed":true,"Failure":{"StatusCode":429}}`)
	for i := 0; i < 3; i++ {
		_ = r.Handle(context.Background(), runtime.MethodUsageHandle, payload)
	}
	if document := readPriorityDocument(t, authPath); strings.Contains(document, `"priority":-1`) {
		t.Fatalf("disabled plugin demoted a credential: %s", document)
	}
	diagnostics, err := r.Diagnostics(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if pending, _ := diagnostics["pending_429"].([]map[string]any); len(pending) != 0 {
		t.Fatalf("disabled plugin still tracked 429s: %#v", pending)
	}
}

func TestRuntime_429CooldownRunsWithAutoApplyOff(t *testing.T) {
	r, authPath := newCooldownCandidate(t, "enabled: true\n")
	payload := []byte(`{"Provider":"antigravity","AuthIndex":"auth","Failed":true,"Failure":{"StatusCode":429}}`)
	_ = r.Handle(context.Background(), runtime.MethodUsageHandle, payload)
	_ = r.Handle(context.Background(), runtime.MethodUsageHandle, payload)
	if document := readPriorityDocument(t, authPath); !strings.Contains(document, `"priority":-1`) {
		t.Fatalf("reactive cooldown should not depend on auto_apply: %s", document)
	}
}

func TestRuntime_FailureMentioning429IsNotARateLimit(t *testing.T) {
	r, authPath := newCooldownCandidate(t, "")
	payload := []byte(`{"Provider":"antigravity","AuthIndex":"auth","Failed":true,"Failure":{"StatusCode":500,"Body":"{\"error\":{\"code\":500,\"message\":\"upstream returned 429 bytes of garbage\",\"status\":\"INTERNAL\"}}"}}`)
	for i := 0; i < 3; i++ {
		_ = r.Handle(context.Background(), runtime.MethodUsageHandle, payload)
	}
	if document := readPriorityDocument(t, authPath); strings.Contains(document, `"priority":-1`) {
		t.Fatalf("a non-rate-limit failure triggered cooldown: %s", document)
	}
}
