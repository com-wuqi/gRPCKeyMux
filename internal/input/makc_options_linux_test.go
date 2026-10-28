//go:build linux

package input

import (
	"testing"

	"GoKeyMux/internal/config"
)

func TestChooseMakcBackend(t *testing.T) {
	cases := []struct {
		name         string
		mode         config.MakcBackendName
		protocol     string
		wantPrimary  config.MakcBackendName
		wantFallback config.MakcBackendName
	}{
		{"explicit_uinput", config.MakcBackendUInput, sessionProtocolWayland, config.MakcBackendUInput, ""},
		{"explicit_portal", config.MakcBackendXDGPortal, sessionProtocolX11, config.MakcBackendXDGPortal, ""},
		{"auto_wayland", config.MakcBackendAuto, sessionProtocolWayland, config.MakcBackendXDGPortal, config.MakcBackendUInput},
		{"auto_x11", config.MakcBackendAuto, sessionProtocolX11, config.MakcBackendUInput, config.MakcBackendXDGPortal},
		{"auto_unknown", config.MakcBackendAuto, sessionProtocolUnknown, config.MakcBackendUInput, config.MakcBackendXDGPortal},
		{"empty_means_auto", "", sessionProtocolWayland, config.MakcBackendXDGPortal, config.MakcBackendUInput},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			primary, fallback := chooseMakcBackend(c.mode, c.protocol)
			if primary != c.wantPrimary || fallback != c.wantFallback {
				t.Fatalf("chooseMakcBackend(%q, %q) = (%q, %q), want (%q, %q)",
					c.mode, c.protocol, primary, fallback, c.wantPrimary, c.wantFallback)
			}
		})
	}
}

func TestDetectSessionProtocol(t *testing.T) {
	cases := []struct {
		name    string
		session string
		wayland string
		display string
		want    string
	}{
		{"xdg_wayland", "wayland", "", ":0", sessionProtocolWayland},
		{"xdg_x11", "x11", "", ":0", sessionProtocolX11},
		{"xdg_xorg", "xorg", "", ":0", sessionProtocolX11},
		{"wayland_display", "tty", "wayland-0", ":0", sessionProtocolWayland},
		{"display_only", "tty", "", ":0", sessionProtocolX11},
		{"unknown", "", "", "", sessionProtocolUnknown},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			t.Setenv("XDG_SESSION_TYPE", c.session)
			t.Setenv("WAYLAND_DISPLAY", c.wayland)
			t.Setenv("DISPLAY", c.display)
			if got := detectSessionProtocol(); got != c.want {
				t.Fatalf("detectSessionProtocol() = %q, want %q", got, c.want)
			}
		})
	}
}
