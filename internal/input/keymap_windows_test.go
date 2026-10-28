//go:build windows

package input

import (
	"testing"

	"github.com/aiwaki/makc"
	"github.com/rpdg/winput"
)

func TestKeyCodesFromWinput(t *testing.T) {
	cases := []struct {
		in   winput.Key
		want KeyCodes
	}{
		{winput.KeyA, KeyCodes{makc.KeyA, fakerInputKeyA, 0}},
		{winput.KeyF1, KeyCodes{makc.KeyF1, fakerInputKeyF1, 0}},
		{winput.KeyLeft, KeyCodes{makc.KeyLeft, fakerInputKeyLeftArrow, 0}},
		{winput.KeyNumLock, KeyCodes{makc.KeyNumLock, fakerInputKeyNumLock, 0}},
	}
	for _, c := range cases {
		got, ok := KeyCodesFromWinput(c.in)
		if !ok || got != c.want {
			t.Errorf("KeyCodesFromWinput(%#x) = (%+v, %v), want (%+v, true)", c.in, got, ok, c.want)
		}
	}
	if _, ok := KeyCodesFromWinput(winput.Key(0xFFFF)); ok {
		t.Errorf("KeyCodesFromWinput(0xFFFF) should not be ok")
	}
}

// TestKeyCodesWinputConsistency verifies that every table entry and every mapped
// rune round-trips through the Windows scan-code map.
func TestKeyCodesWinputConsistency(t *testing.T) {
	if len(keyByWinput) != len(fakerInputKeyTable) {
		t.Fatalf("keyByWinput has %d entries, want %d (duplicate or missing scan codes?)", len(keyByWinput), len(fakerInputKeyTable))
	}

	for _, kc := range fakerInputKeyTable {
		scan, ok := scanByMakc[kc.Makc]
		if !ok {
			t.Errorf("no scan code for makc=%#x", kc.Makc)
			continue
		}
		if got, ok := KeyCodesFromWinput(scan); !ok || got != kc {
			t.Errorf("table entry makc=%#x scan=%#x inconsistent: got (%+v, %v)", kc.Makc, scan, got, ok)
		}
	}

	for r, kc := range keyByRune {
		scan := winputKey(kc)
		if scan == 0 {
			t.Errorf("rune %q has no scan code", r)
			continue
		}
		if got, ok := KeyCodesFromWinput(scan); !ok || got.Makc != kc.Makc || got.FakerInput != kc.FakerInput {
			t.Errorf("rune %q: scan=%#x inconsistent with keyByWinput (got %+v, %v)", r, scan, got, ok)
		}
	}
}

// TestFakerInputKeyFromWinput verifies the FakerInput-only scan-code helper.
func TestFakerInputKeyFromWinput(t *testing.T) {
	code, ok := FakerInputKeyFromWinput(winput.KeyEsc)
	if !ok || code != fakerInputKeyEscape {
		t.Errorf("FakerInputKeyFromWinput(Esc) = (%#x, %v), want (%#x, true)", code, ok, fakerInputKeyEscape)
	}
}

// TestWinputKeyModifiers covers mappings that are not part of the shared key
// table (standalone modifiers) plus a key winput cannot represent.
func TestWinputKeyModifiers(t *testing.T) {
	cases := []struct {
		name string
		want winput.Key
	}{
		{"ctrl", winput.KeyCtrl},
		{"lctrl", winput.KeyCtrl},
		{"rctrl", winput.KeyCtrl},
		{"shift", winput.KeyShift},
		{"alt", winput.KeyAlt},
		{"gui", 0},
	}
	for _, c := range cases {
		kc, ok := KeyCodesFromName(c.name)
		if !ok {
			t.Fatalf("KeyCodesFromName(%q) failed", c.name)
		}
		if got := winputKey(kc); got != c.want {
			t.Errorf("winputKey(%q) = %#x, want %#x", c.name, got, c.want)
		}
	}
}
