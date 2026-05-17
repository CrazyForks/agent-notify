package agentintegrations

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/hellolib/agent-notify/internal/codexhooks"
	"github.com/hellolib/agent-notify/internal/common"
)

// CodexIntegration implements Integration for Codex.
type CodexIntegration struct{}

// NewCodexIntegration creates a new Codex integration.
func NewCodexIntegration() *CodexIntegration {
	return &CodexIntegration{}
}

// Name returns the display name for Codex.
func (c *CodexIntegration) Name() string {
	return "Codex"
}

// DetectInstalled checks if the codex CLI is installed.
func (c *CodexIntegration) DetectInstalled() bool {
	_, err := exec.LookPath("codex")
	return err == nil
}

// SettingsPath returns the path to Codex's hooks.json file.
// Codex 同时支持 ~/.codex/hooks.json 与 ~/.codex/config.toml 内联 [hooks]；
// 这里统一使用 hooks.json，结构上与 Claude settings.json 对齐，便于维护。
func (c *CodexIntegration) SettingsPath(scope string) (string, error) {
	switch scope {
	case "user":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".codex", "hooks.json"), nil
	case "project":
		return filepath.Join(".codex", "hooks.json"), nil
	default:
		return "", fmt.Errorf("unsupported scope: %s", scope)
	}
}

// Install 写入 Codex hooks.json，订阅 PermissionRequest 与 Stop 事件。
func (c *CodexIntegration) Install(settingsPath, binaryPath string) error {
	return codexhooks.Install(settingsPath, common.ResolveBinaryPath(binaryPath))
}

// IsHookInstalled 检查 Codex hooks.json 中是否已经登记了 handle-codex-hook。
func (c *CodexIntegration) IsHookInstalled(settingsPath string) (bool, error) {
	data, err := os.ReadFile(settingsPath)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	settings := map[string]any{}
	if err := json.Unmarshal(data, &settings); err != nil {
		return false, fmt.Errorf("failed to parse hooks.json: %w", err)
	}

	hooks, ok := settings["hooks"].(map[string]any)
	if !ok {
		return false, nil
	}

	pr, ok := hooks["PermissionRequest"].([]any)
	if !ok || len(pr) == 0 {
		return false, nil
	}

	entry, ok := pr[0].(map[string]any)
	if !ok {
		return false, nil
	}

	hookList, ok := entry["hooks"].([]any)
	if !ok || len(hookList) == 0 {
		return false, nil
	}

	hook, ok := hookList[0].(map[string]any)
	if !ok {
		return false, nil
	}

	cmd, ok := hook["command"].(string)
	if !ok {
		return false, nil
	}

	return strings.Contains(cmd, "handle-codex-hook"), nil
}
