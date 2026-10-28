//go:build windows

package input

import (
	"fmt"
	"log/slog"
	"time"

	"GoKeyMux/internal/config"
)

// NewBackend constructs the backend selected by driveName. All backends are
// available on Windows.
func NewBackend(driveName config.DriveName, noopLatency time.Duration) (Backend, error) {
	switch driveName {
	case config.DriveMakc:
		return newMakcBackend()
	case config.DriveFakerInput:
		return newFakerInputBackend()
	case config.DriveWinputWithWindow:
		slog.Warn("使用winput可能会导致键鼠输入串行化")
		return newWinputWindowBackend()
	case config.DriveWinputWithInterception:
		slog.Warn("使用winput可能会导致键鼠输入串行化")
		return newWinputInterceptionBackend()
	case config.DriveNoop:
		slog.Debug("debug driver ready")
		return NewNoopDevice(noopLatency), nil
	default:
		return nil, fmt.Errorf("unsupported drive %s", driveName)
	}
}
