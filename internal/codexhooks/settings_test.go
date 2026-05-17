package codexhooks

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestBuildHookSettings_RegistersTwoEvents(t *testing.T) {
	got := BuildHookSettings("/tmp/agent-notify")

	hooks, ok := got["hooks"].(map[string]any)
	if !ok {
		t.Fatalf("hooks type = %T, want map[string]any", got["hooks"])
	}

	for _, event := range []string{"PermissionRequest", "Stop"} {
		items, ok := hooks[event].([]map[string]any)
		if !ok || len(items) != 1 {
			t.Fatalf("%s entries missing or invalid: %v", event, hooks[event])
		}
		entryHooks, ok := items[0]["hooks"].([]map[string]any)
		if !ok || len(entryHooks) != 1 {
			t.Fatalf("%s command list missing or invalid", event)
		}
		if entryHooks[0]["command"] != "/tmp/agent-notify handle-codex-hook" {
			t.Fatalf("%s command = %v, want /tmp/agent-notify handle-codex-hook", event, entryHooks[0]["command"])
		}
		if entryHooks[0]["type"] != "command" {
			t.Fatalf("%s type = %v, want command", event, entryHooks[0]["type"])
		}
	}

	// 不应注册 Codex 不支持的事件
	for _, unsupported := range []string{"Notification", "PostToolUseFailure", "UserPromptSubmit", "PreToolUse", "PostToolUse", "SessionStart"} {
		if _, exists := hooks[unsupported]; exists {
			t.Fatalf("hooks should not contain %s for Codex", unsupported)
		}
	}
}

func TestInstall_MergesExistingHooks(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "hooks.json")
	existing := `{"hooks":{"SessionStart":[{"hooks":[{"type":"command","command":"echo hi"}]}]}}`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := Install(path, "/tmp/agent-notify"); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	var got map[string]any
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	hooks, ok := got["hooks"].(map[string]any)
	if !ok {
		t.Fatal("hooks key missing or wrong type")
	}
	for _, key := range []string{"SessionStart", "PermissionRequest", "Stop"} {
		if _, exists := hooks[key]; !exists {
			t.Fatalf("hooks missing key %q after install", key)
		}
	}
}

func TestInstall_CreatesParentDirectories(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "nested", "deeper", "hooks.json")

	if err := Install(path, "/tmp/agent-notify"); err != nil {
		t.Fatalf("Install() error = %v", err)
	}

	if _, err := os.Stat(path); err != nil {
		t.Fatalf("hooks.json not created at %q: %v", path, err)
	}
}
