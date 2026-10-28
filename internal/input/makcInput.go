package input

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/aiwaki/makc"
)

// makcScreenSize caches the primary screen size captured at startup. It is used
// to scale normalized 0..10000 absolute mouse coordinates into pixels. It is
// written by MakcInputInit before the engine workers start and, if the startup
// query failed (e.g. Wayland), by the mouse worker on the first successful
// query. After startup only the mouse worker touches it, so no further
// synchronization is required.
var makcScreenSize makc.Point

func MakcInputInit() (*makc.Client, error) {
	slog.Debug("MakcInputInit")

	plan := makcOpenPlanForPlatform()
	client, err := makc.Open(plan.primary...)
	if err != nil && len(plan.fallback) > 0 {
		slog.Warn("makc primary backend failed, trying fallback", "err", err)
		client, err = makc.Open(plan.fallback...)
	}
	if err != nil {
		return nil, err
	}
	logMakcBackends(client)

	if size, err := client.Mouse.ScreenSize(context.Background()); err == nil {
		makcScreenSize = size
		slog.Debug("Makc ScreenSize: ", "X", size.X, "Y", size.Y)
	} else {
		// Wayland has no X11 display, so absolute movement is unavailable; the
		// engine reports that explicitly when ABS_MOVE is requested.
		slog.Warn("MakcInputInit: query screen size failed (absolute mouse movement unavailable)", "err", err)
	}
	return client, nil
}

func MakcInputPress(client *makc.Client, key KeyCodes, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputPress failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	return client.Keyboard.Press(ctx, key.Makc)
}

func MakcInputRelease(client *makc.Client, key KeyCodes, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputRelease failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	return client.Keyboard.Release(ctx, key.Makc)
}

// MakcMouseButton maps a backend-neutral MouseButton to a makc.MouseButton.
func MakcMouseButton(b MouseButton) (makc.MouseButton, error) {
	switch b {
	case MouseButtonLeft:
		return makc.ButtonLeft, nil
	case MouseButtonRight:
		return makc.ButtonRight, nil
	case MouseButtonMiddle:
		return makc.ButtonMiddle, nil
	case MouseButtonX1:
		return makc.ButtonX1, nil
	case MouseButtonX2:
		return makc.ButtonX2, nil
	default:
		return 0, fmt.Errorf("unknown mouse button %d", b)
	}
}

// MakcAbsCoordinate scales a normalized 0..10000 coordinate to screen pixels,
// mapping onto the inclusive pixel range [0, screenSize-1] so that 10000 lands
// on the last valid pixel (screenSize-1). This mirrors FakerInput's absolute
// range, where its logical maximum (32767) equals the screen's bottom/right edge.
func MakcAbsCoordinate(v int32, screenSize int) int {
	if v < 0 {
		v = 0
	}
	if v > 10000 {
		v = 10000
	}
	if screenSize <= 1 {
		return 0
	}
	return int(int64(v) * int64(screenSize-1) / 10000)
}

func MakcInputMouseMove(client *makc.Client, dx, dy int32, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputMouseMove failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	return client.Mouse.MoveBy(ctx, int(dx), int(dy))
}

func MakcInputMouseMoveTo(client *makc.Client, x, y int32, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputMouseMoveTo failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	size := makcScreenSize
	if size.X <= 0 || size.Y <= 0 {
		var err error
		size, err = client.Mouse.ScreenSize(ctx)
		if err != nil {
			return fmt.Errorf("absolute mouse movement is unavailable: screen size query failed: %w", err)
		}
		makcScreenSize = size
	}
	return client.Mouse.MoveTo(ctx, MakcAbsCoordinate(x, size.X), MakcAbsCoordinate(y, size.Y))
}

func MakcInputMouseButtonDown(client *makc.Client, button MouseButton, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputMouseButtonDown failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	b, err := MakcMouseButton(button)
	if err != nil {
		return err
	}
	return client.Mouse.Press(ctx, b)
}

func MakcInputMouseButtonUp(client *makc.Client, button MouseButton, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputMouseButtonUp failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	b, err := MakcMouseButton(button)
	if err != nil {
		return err
	}
	return client.Mouse.Release(ctx, b)
}

func MakcInputMouseWheel(client *makc.Client, delta int32, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputMouseWheel failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	return client.Mouse.Wheel(ctx, int(delta))
}

func MakcInputMouseHWheel(client *makc.Client, delta int32, ctx context.Context) error {
	if client == nil {
		slog.Error("MakcInputMouseHWheel failed: client is nil")
		return fmt.Errorf("makc client is nil")
	}
	return client.Mouse.HWheel(ctx, int(delta))
}

// MakcScreenSize returns the screen size cached by MakcInputInit. It is used by
// integration tests; production code reads the cache directly.
func MakcScreenSize() (width, height int) {
	return makcScreenSize.X, makcScreenSize.Y
}
