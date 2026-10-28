//go:build windows

package input

import (
	"context"
	"fmt"
)

type fakerInputBackend struct {
	client *FakerInputDevice
}

func newFakerInputBackend() (Backend, error) {
	client, err := FakerInputInit()
	if err != nil {
		return nil, err
	}
	return &fakerInputBackend{client: client}, nil
}

func (b *fakerInputBackend) Name() string { return "fakerInput" }

func (b *fakerInputBackend) Press(_ context.Context, key KeyCodes) error {
	return b.client.KeyDown(key.FakerInput, key.Modifiers)
}

func (b *fakerInputBackend) Release(_ context.Context, key KeyCodes) error {
	return b.client.KeyUp(key.FakerInput, key.Modifiers)
}

func (b *fakerInputBackend) Mouse(_ context.Context, ev MouseEvent) error {
	switch ev.Kind {
	case MouseEventMoveRel:
		return b.client.MouseMove(FakerInputRelDelta(ev.DX), FakerInputRelDelta(ev.DY))
	case MouseEventMoveAbs:
		return b.client.MouseMoveTo(FakerInputAbsCoordinate(ev.X), FakerInputAbsCoordinate(ev.Y))
	case MouseEventButton:
		button, err := FakerInputMouseButton(ev.Button)
		if err != nil {
			return err
		}
		if ev.IsPressed {
			return b.client.MouseButtonDown(button)
		}
		return b.client.MouseButtonUp(button)
	case MouseEventWheel:
		return b.client.MouseWheel(int8(ev.Delta))
	case MouseEventHWheel:
		return b.client.MouseHWheel(int8(ev.Delta))
	default:
		return fmt.Errorf("unknown mouse event %d", ev.Kind)
	}
}

func (b *fakerInputBackend) Close() error {
	if b == nil || b.client == nil {
		return nil
	}
	return b.client.Close()
}

var _ Backend = (*fakerInputBackend)(nil)
