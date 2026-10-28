//go:build linux

package input

import (
	"log/slog"
	"os"
	"strings"

	"GoKeyMux/internal/config"

	"github.com/aiwaki/makc"
)

const (
	sessionProtocolWayland = "wayland"
	sessionProtocolX11     = "x11"
	sessionProtocolUnknown = "unknown"
)

// detectSessionProtocol infers the display protocol. XDG_SESSION_TYPE is the
// authoritative source when present; WAYLAND_DISPLAY/DISPLAY cover services
// started without it (for example by systemd).
func detectSessionProtocol() string {
	switch strings.ToLower(strings.TrimSpace(os.Getenv("XDG_SESSION_TYPE"))) {
	case "wayland":
		return sessionProtocolWayland
	case "x11", "xorg":
		return sessionProtocolX11
	}
	if os.Getenv("WAYLAND_DISPLAY") != "" {
		return sessionProtocolWayland
	}
	if os.Getenv("DISPLAY") != "" {
		return sessionProtocolX11
	}
	return sessionProtocolUnknown
}

// chooseMakcBackend resolves the configured mode for one input class. An
// explicit value is used as-is; auto prefers the protocol's native backend and
// offers the other one as a fallback.
func chooseMakcBackend(mode config.MakcBackendName, protocol string) (primary, fallback config.MakcBackendName) {
	switch mode {
	case config.MakcBackendUInput, config.MakcBackendXDGPortal:
		return mode, ""
	}
	if protocol == sessionProtocolWayland {
		return config.MakcBackendXDGPortal, config.MakcBackendUInput
	}
	return config.MakcBackendUInput, config.MakcBackendXDGPortal
}

func makcOptions(keyboardBackend, mouseBackend config.MakcBackendName) []makc.Option {
	opts := make([]makc.Option, 0, 2)
	switch keyboardBackend {
	case config.MakcBackendXDGPortal:
		opts = append(opts, makc.WithKeyboardInjection(makc.KeyboardInjectionXDGPortal))
	default:
		opts = append(opts, makc.WithKeyboardInjection(makc.KeyboardInjectionUInput))
	}
	switch mouseBackend {
	case config.MakcBackendXDGPortal:
		opts = append(opts, makc.WithMouseInjection(makc.MouseInjectionXDGPortal))
	default:
		opts = append(opts, makc.WithMouseInjection(makc.MouseInjectionUInput))
	}
	return opts
}

func makcOpenPlanForPlatform() makcOpenPlan {
	protocol := detectSessionProtocol()

	kb, kbFallback := chooseMakcBackend(config.GlobalConfig.MakcKeyboardBackend, protocol)
	ms, msFallback := chooseMakcBackend(config.GlobalConfig.MakcMouseBackend, protocol)

	slog.Info("makc backend selection",
		"protocol", protocol,
		"keyboard", string(kb),
		"mouse", string(ms),
	)

	plan := makcOpenPlan{primary: makcOptions(kb, ms)}

	if kbFallback == "" {
		kbFallback = kb
	}
	if msFallback == "" {
		msFallback = ms
	}
	if kbFallback != kb || msFallback != ms {
		plan.fallback = makcOptions(kbFallback, msFallback)
		slog.Info("makc fallback backend available",
			"keyboard", string(kbFallback),
			"mouse", string(msFallback),
		)
	}
	return plan
}
