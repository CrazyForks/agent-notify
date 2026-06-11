package cli

import (
	"context"
	"io"

	"github.com/hellolib/agent-notify/internal/config"
	"github.com/hellolib/agent-notify/internal/i18n"
)

func Run(ctx context.Context, args []string, stdin io.Reader, stdout, stderr io.Writer) error {
	initLocale()

	streams := Streams{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	}
	if len(args) == 0 {
		return runMenu(ctx, streams)
	}

	cmd := NewRootCmd(ctx, Streams{
		Stdin:  stdin,
		Stdout: stdout,
		Stderr: stderr,
	})
	cmd.SetArgs(args)
	return cmd.Execute()
}

// initLocale loads the persisted locale from config and applies it.
// If the config cannot be loaded, the default (zh-CN) is used.
func initLocale() {
	path, err := config.DefaultPath()
	if err != nil {
		return
	}
	cfg, err := config.Load(path)
	if err != nil {
		return
	}
	i18n.Set(cfg.Behavior.Locale)
}
