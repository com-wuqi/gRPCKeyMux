//go:build !windows

package input

import (
	"fmt"
	"log/slog"
	"time"

	"GoKeyMux/internal/config"
)

// NewBackend constructs the backend selected by driveName. Only makc and noop
// are available outside Windows; the Windows-only drives fail with an explicit
// error rather than a compile-time break.
func NewBackend(driveName config.DriveName, noopLatency time.Duration) (Backend, error) {
	switch driveName {
	case config.DriveMakc:
		return newMakcBackend()
	case config.DriveNoop:
		slog.Debug("debug driver ready")
		return NewNoopDevice(noopLatency), nil
	case config.DriveFakerInput, config.DriveWinputWithWindow, config.DriveWinputWithInterception:
		return nil, fmt.Errorf("drive %q is only supported on Windows; use \"makc\" or \"noop\"", driveName)
	default:
		return nil, fmt.Errorf("unsupported drive %s", driveName)
	}
}
