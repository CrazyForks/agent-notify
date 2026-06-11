// Package i18n provides internationalization support for the interactive TUI.
// The hook notification text (internal/notify/format.go) is not affected.
package i18n

import "sync"

// Lang represents a supported language.
type Lang string

const (
	ZhCN Lang = "zh-CN"
	EnUS Lang = "en-US"
)

var (
	mu      sync.RWMutex
	current Lang = ZhCN
)

// Set sets the current language. Any locale that is not "en-US" is normalized to "zh-CN".
func Set(locale string) {
	mu.Lock()
	defer mu.Unlock()
	if locale == "en-US" {
		current = EnUS
	} else {
		current = ZhCN
	}
}

// Current returns the current language code.
func Current() string {
	mu.RLock()
	defer mu.RUnlock()
	return string(current)
}

// IsEnglish reports whether the current language is English.
func IsEnglish() bool {
	mu.RLock()
	defer mu.RUnlock()
	return current == EnUS
}

// T returns the translated string for the given key.
// It falls back to the Chinese string if the key is missing for the current language,
// and falls back to the key itself if missing for both languages.
func T(key string) string {
	mu.RLock()
	lang := current
	mu.RUnlock()

	langs, ok := catalog[key]
	if !ok {
		return key
	}
	if s, ok := langs[lang]; ok {
		return s
	}
	if s, ok := langs[ZhCN]; ok {
		return s
	}
	return key
}
