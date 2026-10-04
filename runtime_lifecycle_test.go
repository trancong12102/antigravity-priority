package main

import (
	"encoding/json"
	"path/filepath"
	"testing"

	pluginruntime "antigravity-priority/internal/runtime"
)

func TestRuntimeIsNotCreatedWhenTheLibraryLoads(t *testing.T) {
	if cpaRuntime.Load() != nil {
		t.Fatal("package initialisation created a Runtime before cliproxy_plugin_init")
	}
	var envelope struct {
		OK    bool `json:"ok"`
		Error struct {
			Code string `json:"code"`
		} `json:"error"`
	}
	if err := json.Unmarshal(handlePluginCall("plugin.register", []byte(`{}`)), &envelope); err != nil {
		t.Fatal(err)
	}
	if envelope.OK || envelope.Error.Code != "not_initialized" {
		t.Fatalf("call before init should fail cleanly, got %+v", envelope)
	}
}

func TestRepeatedInitShutsDownThePreviousRuntime(t *testing.T) {
	t.Cleanup(func() {
		if rt := cpaRuntime.Swap(nil); rt != nil {
			_ = rt.Shutdown(t.Context())
		}
	})
	options := pluginruntime.Options{Host: hostCallbackAdapter{}, StateCachePath: filepath.Join(t.TempDir(), "cache.json")}
	initRuntime(options)
	first := cpaRuntime.Load()
	initRuntime(options)
	second := cpaRuntime.Load()
	if first == nil || second == nil || first == second {
		t.Fatalf("expected two distinct runtimes, got %p and %p", first, second)
	}
	register := []byte(`{"config_yaml":"enabled: true\n"}`)
	if envelopeOK(t, first.Handle(t.Context(), pluginruntime.MethodPluginRegister, register)) {
		t.Fatal("previous runtime still accepts registration; it was not shut down")
	}
	if !envelopeOK(t, second.Handle(t.Context(), pluginruntime.MethodPluginRegister, register)) {
		t.Fatal("current runtime rejected registration")
	}
}

func envelopeOK(t *testing.T, raw []byte) bool {
	t.Helper()
	var envelope struct {
		OK bool `json:"ok"`
	}
	if err := json.Unmarshal(raw, &envelope); err != nil {
		t.Fatal(err)
	}
	return envelope.OK
}
