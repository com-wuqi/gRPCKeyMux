package input

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aiwaki/makc"
)

// MakcBackend adapts *makc.Client to the Backend interface. It is available on
// every platform supported by makc (Windows, Linux, macOS).
type MakcBackend struct {
	Client *makc.Client
}

func newMakcBackend() (Backend, error) {
	client, err := MakcInputInit()
	if err != nil {
		return nil, err
	}
	return &MakcBackend{Client: client}, nil
}

func (b *MakcBackend) Name() string { return "makc" }

func (b *MakcBackend) Press(ctx context.Context, key KeyCodes) error {
	return MakcInputPress(b.Client, key, ctx)
}

func (b *MakcBackend) Release(ctx context.Context, key KeyCodes) error {
	return MakcInputRelease(b.Client, key, ctx)
}

func (b *MakcBackend) Mouse(ctx context.Context, ev MouseEvent) error {
	switch ev.Kind {
	case MouseEventMoveRel:
		return MakcInputMouseMove(b.Client, ev.DX, ev.DY, ctx)
	case MouseEventMoveAbs:
		return MakcInputMouseMoveTo(b.Client, ev.X, ev.Y, ctx)
	case MouseEventButton:
		if ev.IsPressed {
			return MakcInputMouseButtonDown(b.Client, ev.Button, ctx)
		}
		return MakcInputMouseButtonUp(b.Client, ev.Button, ctx)
	case MouseEventWheel:
		return MakcInputMouseWheel(b.Client, ev.Delta, ctx)
	case MouseEventHWheel:
		return MakcInputMouseHWheel(b.Client, ev.Delta, ctx)
	default:
		return fmt.Errorf("unknown mouse event %d", ev.Kind)
	}
}

func (b *MakcBackend) Close() error {
	if b == nil || b.Client == nil {
		return nil
	}
	return b.Client.Close()
}

// ensure the adapter satisfies the interface at compile time.
var _ Backend = (*MakcBackend)(nil)

// logMakcBackends is called after a client is opened so the selected
// platform-specific injection backends are visible in the startup log.
func logMakcBackends(client *makc.Client) {
	if client == nil {
		return
	}
	slog.Info("makc injection backends",
		"keyboard", client.Keyboard.InjectionBackend().String(),
		"mouse", client.Mouse.InjectionBackend().String(),
	)
}
