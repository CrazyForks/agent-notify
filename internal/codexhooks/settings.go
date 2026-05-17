package codexhooks

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"path/filepath"

	"github.com/hellolib/agent-notify/internal/common"
)

// BuildHookSettings 生成 Codex hooks.json 所需的 settings 结构。
// Codex 当前可靠支持的事件只有 PermissionRequest 与 Stop，
// 分别对应项目里的 permission_required / run_completed。
func BuildHookSettings(binaryPath string) map[string]any {
	binaryPath = common.ResolveBinaryPath(binaryPath)
	command := binaryPath + " handle-codex-hook"

	buildEntry := func() []map[string]any {
		return []map[string]any{
			{
				"hooks": []map[string]any{
					{
						"type":    "command",
						"command": command,
					},
				},
			},
		}
	}

	return map[string]any{
		"hooks": map[string]any{
			"PermissionRequest": buildEntry(),
			"Stop":              buildEntry(),
		},
	}
}

func Install(path string, binaryPath string) error {
	settings := map[string]any{}

	data, err := os.ReadFile(path)
	if err == nil {
		if err := json.Unmarshal(data, &settings); err != nil {
			return err
		}
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}

	builtHooks := BuildHookSettings(binaryPath)["hooks"].(map[string]any)
	existingHooks, _ := settings["hooks"].(map[string]any)
	if existingHooks == nil {
		existingHooks = map[string]any{}
	}
	for key, value := range builtHooks {
		existingHooks[key] = value
	}
	settings["hooks"] = existingHooks

	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(path, out, 0o644); err != nil {
		return err
	}

	// 同步启用 config.toml 中的 [features] hooks = true
	configTomlPath := filepath.Join(filepath.Dir(path), "config.toml")
	if err := EnableHooksFeature(configTomlPath); err != nil {
		log.Printf("warning: failed to enable hooks feature in %s: %v", configTomlPath, err)
	}

	return nil
}
