//go:build windows

package input

import (
	"context"
	"fmt"

	"github.com/rpdg/winput"
)

type winputWindowBackend struct {
	window *winput.Window
}

func newWinputWindowBackend() (Backend, error) {
	window, err := WinputInitWithWindow()
	if err != nil {
		return nil, err
	}
	return &winputWindowBackend{window: window}, nil
}

func (b *winputWindowBackend) Name() string { return "winputWithWindow" }

func (b *winputWindowBackend) Press(_ context.Context, key KeyCodes) error {
	return WinputWithWindowPress(b.window, key)
}

func (b *winputWindowBackend) Release(_ context.Context, key KeyCodes) error {
	return WinputWithWindowRelease(b.window, key)
}

func (b *winputWindowBackend) Mouse(_ context.Context, _ MouseEvent) error {
	return fmt.Errorf("mouse not implemented for drive %s", b.Name())
}

func (b *winputWindowBackend) Close() error { return nil }

type winputInterceptionBackend struct{}

func newWinputInterceptionBackend() (Backend, error) {
	if err := WinputWithInterceptionInit(); err != nil {
		return nil, err
	}
	return &winputInterceptionBackend{}, nil
}

func (b *winputInterceptionBackend) Name() string { return "winputWithInterception" }

func (b *winputInterceptionBackend) Press(_ context.Context, key KeyCodes) error {
	return WinputWithInterceptionPress(key)
}

func (b *winputInterceptionBackend) Release(_ context.Context, key KeyCodes) error {
	return WinputWithInterceptionRelease(key)
}

func (b *winputInterceptionBackend) Mouse(_ context.Context, _ MouseEvent) error {
	return fmt.Errorf("mouse not implemented for drive %s", b.Name())
}

func (b *winputInterceptionBackend) Close() error { return nil }

var (
	_ Backend = (*winputWindowBackend)(nil)
	_ Backend = (*winputInterceptionBackend)(nil)
)
