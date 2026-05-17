package codexhooks

import (
	"os"
	"path/filepath"
	"testing"

	toml "github.com/pelletier/go-toml/v2"
)

func TestEnableHooksFeature_CreatesNewFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	if err := EnableHooksFeature(path); err != nil {
		t.Fatalf("EnableHooksFeature() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config.toml: %v", err)
	}

	config := map[string]any{}
	if err := toml.Unmarshal(data, &config); err != nil {
		t.Fatalf("failed to parse config.toml: %v", err)
	}

	features, ok := config["features"].(map[string]any)
	if !ok {
		t.Fatal("features section missing or wrong type")
	}
	if features["hooks"] != true {
		t.Fatalf("features.hooks = %v, want true", features["hooks"])
	}
}

func TestEnableHooksFeature_AppendsToExisting(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	existing := `model = "gpt-4"
`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := EnableHooksFeature(path); err != nil {
		t.Fatalf("EnableHooksFeature() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config.toml: %v", err)
	}

	config := map[string]any{}
	if err := toml.Unmarshal(data, &config); err != nil {
		t.Fatalf("failed to parse config.toml: %v", err)
	}

	if config["model"] != "gpt-4" {
		t.Fatalf("model = %v, want gpt-4 (existing config should be preserved)", config["model"])
	}

	features, ok := config["features"].(map[string]any)
	if !ok {
		t.Fatal("features section missing or wrong type")
	}
	if features["hooks"] != true {
		t.Fatalf("features.hooks = %v, want true", features["hooks"])
	}
}

func TestEnableHooksFeature_AddsToExistingFeatures(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	existing := `[features]
fast_mode = true
`
	if err := os.WriteFile(path, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	if err := EnableHooksFeature(path); err != nil {
		t.Fatalf("EnableHooksFeature() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config.toml: %v", err)
	}

	config := map[string]any{}
	if err := toml.Unmarshal(data, &config); err != nil {
		t.Fatalf("failed to parse config.toml: %v", err)
	}

	features, ok := config["features"].(map[string]any)
	if !ok {
		t.Fatal("features section missing or wrong type")
	}
	if features["hooks"] != true {
		t.Fatalf("features.hooks = %v, want true", features["hooks"])
	}
	if features["fast_mode"] != true {
		t.Fatalf("features.fast_mode = %v, want true (existing feature should be preserved)", features["fast_mode"])
	}
}

func TestEnableHooksFeature_Idempotent(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")

	if err := EnableHooksFeature(path); err != nil {
		t.Fatalf("first call error = %v", err)
	}
	if err := EnableHooksFeature(path); err != nil {
		t.Fatalf("second call error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read config.toml: %v", err)
	}

	config := map[string]any{}
	if err := toml.Unmarshal(data, &config); err != nil {
		t.Fatalf("failed to parse config.toml: %v", err)
	}

	features, ok := config["features"].(map[string]any)
	if !ok {
		t.Fatal("features section missing or wrong type")
	}
	if features["hooks"] != true {
		t.Fatalf("features.hooks = %v, want true", features["hooks"])
	}
}
