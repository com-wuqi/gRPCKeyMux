package input

import (
	"testing"

	"github.com/aiwaki/makc"
)

// TestKeyCodesFromRune verifies a representative sample of runes against
// independently known values for the portable encodings. Values are Windows VK
// codes and USB HID usage codes respectively.
func TestKeyCodesFromRune(t *testing.T) {
	cases := []struct {
		r    rune
		want KeyCodes
	}{
		{'a', KeyCodes{0x41, 0x04, 0}},
		{'A', KeyCodes{0x41, 0x04, ModLShift}},
		{'z', KeyCodes{0x5A, 0x1D, 0}},
		{'Z', KeyCodes{0x5A, 0x1D, ModLShift}},
		{'0', KeyCodes{0x30, 0x27, 0}},
		{'9', KeyCodes{0x39, 0x26, 0}},
		{'1', KeyCodes{0x31, 0x1E, 0}},
		{'!', KeyCodes{0x31, 0x1E, ModLShift}},
		{' ', KeyCodes{0x20, 0x2C, 0}},
		{'\n', KeyCodes{0x0D, 0x28, 0}},
		{'\r', KeyCodes{0x0D, 0x28, 0}},
		{'\t', KeyCodes{0x09, 0x2B, 0}},
		{'\b', KeyCodes{0x08, 0x2A, 0}},
		{'\x1b', KeyCodes{0x1B, 0x29, 0}},
		{'-', KeyCodes{0xBD, 0x2D, 0}},
		{'_', KeyCodes{0xBD, 0x2D, ModLShift}},
		{'/', KeyCodes{0xBF, 0x38, 0}},
		{'?', KeyCodes{0xBF, 0x38, ModLShift}},
	}
	for _, c := range cases {
		got, ok := KeyCodesFromRune(c.r)
		if !ok {
			t.Errorf("KeyCodesFromRune(%q): unexpected not-ok", c.r)
			continue
		}
		if got != c.want {
			t.Errorf("KeyCodesFromRune(%q) = %+v, want %+v", c.r, got, c.want)
		}
	}

	if _, ok := KeyCodesFromRune('é'); ok {
		t.Errorf("KeyCodesFromRune('é') should not be ok")
	}
}

func TestKeyCodesFromMakc(t *testing.T) {
	cases := []struct {
		in   makc.Key
		want KeyCodes
	}{
		{makc.KeyA, KeyCodes{makc.KeyA, fakerInputKeyA, 0}},
		{makc.KeyF1, KeyCodes{makc.KeyF1, fakerInputKeyF1, 0}},
		{makc.KeyLeft, KeyCodes{makc.KeyLeft, fakerInputKeyLeftArrow, 0}},
		{makc.KeyNumLock, KeyCodes{makc.KeyNumLock, fakerInputKeyNumLock, 0}},
	}
	for _, c := range cases {
		got, ok := KeyCodesFromMakc(c.in)
		if !ok || got != c.want {
			t.Errorf("KeyCodesFromMakc(%#x) = (%+v, %v), want (%+v, true)", c.in, got, ok, c.want)
		}
	}
	if _, ok := KeyCodesFromMakc(makc.Key(0xFFFF)); ok {
		t.Errorf("KeyCodesFromMakc(0xFFFF) should not be ok")
	}
}

// TestKeyCodesConsistency cross-checks every table entry against the portable
// encodings, guaranteeing that the rune table and the makc index agree on the
// same logical key.
func TestKeyCodesConsistency(t *testing.T) {
	if len(keyByMakc) != len(fakerInputKeyTable) {
		t.Fatalf("keyByMakc has %d entries, want %d (duplicate VK codes?)", len(keyByMakc), len(fakerInputKeyTable))
	}

	for _, kc := range fakerInputKeyTable {
		if got, ok := KeyCodesFromMakc(kc.Makc); !ok || got != kc {
			t.Errorf("table entry makc=%#x inconsistent with keyByMakc: got (%+v, %v)", kc.Makc, got, ok)
		}
	}

	for r, kc := range keyByRune {
		if got, ok := KeyCodesFromMakc(kc.Makc); !ok || got.Makc != kc.Makc || got.FakerInput != kc.FakerInput {
			t.Errorf("rune %q: makc=%#x inconsistent with keyByMakc (got %+v, %v)", r, kc.Makc, got, ok)
		}
	}
}

// TestKeyTableReference verifies every mapped rune against an independently
// hard-coded reference of Windows VK codes and USB HID usage codes. This catches
// wrong pairings or off-by-one HID values that the self-consistency test cannot.
func TestKeyTableReference(t *testing.T) {
	check := func(r rune, vk, hid, mod uint32) {
		t.Helper()
		want := KeyCodes{makc.Key(vk), byte(hid), byte(mod)}
		if got, ok := KeyCodesFromRune(r); !ok || got != want {
			t.Errorf("KeyCodesFromRune(%q) = (%+v, %v), want %+v", r, got, ok, want)
		}
	}

	// Letters: VK 0x41+i, HID 0x04+i.
	for i := range 26 {
		vk := uint32(0x41 + i)
		hid := uint32(0x04 + i)
		check(rune('a'+i), vk, hid, 0)
		check(rune('A'+i), vk, hid, uint32(ModLShift))
	}

	// Digits 1-9: VK 0x31+i, HID 0x1E+i. '0' is the exception.
	shiftedDigits := []rune{')', '!', '@', '#', '$', '%', '^', '&', '*', '('}
	for i := 1; i <= 9; i++ {
		vk := uint32(0x31 + i - 1)
		hid := uint32(0x1E + i - 1)
		check(rune('0'+i), vk, hid, 0)
		check(shiftedDigits[i], vk, hid, uint32(ModLShift))
	}
	check('0', 0x30, 0x27, 0)
	check(')', 0x30, 0x27, uint32(ModLShift))

	// Symbols.
	symbols := []struct {
		r   rune
		vk  uint32
		hid uint32
		mod uint32
	}{
		{'`', 0xC0, 0x35, 0},
		{'-', 0xBD, 0x2D, 0},
		{'=', 0xBB, 0x2E, 0},
		{'[', 0xDB, 0x2F, 0},
		{']', 0xDD, 0x30, 0},
		{'\\', 0xDC, 0x31, 0},
		{';', 0xBA, 0x33, 0},
		{'\'', 0xDE, 0x34, 0},
		{',', 0xBC, 0x36, 0},
		{'.', 0xBE, 0x37, 0},
		{'/', 0xBF, 0x38, 0},
	}
	shifted := []rune{'~', '_', '+', '{', '}', '|', ':', '"', '<', '>', '?'}
	for i, s := range symbols {
		check(s.r, s.vk, s.hid, s.mod)
		check(shifted[i], s.vk, s.hid, uint32(ModLShift))
	}

	// Control runes.
	check(' ', 0x20, 0x2C, 0)
	check('\n', 0x0D, 0x28, 0)
	check('\r', 0x0D, 0x28, 0)
	check('\t', 0x09, 0x2B, 0)
	check('\b', 0x08, 0x2A, 0)
	check('\x1b', 0x1B, 0x29, 0)

	// Every mapped rune must be covered by the reference above; otherwise the
	// reference check is incomplete.
	if len(keyByRune) != 26*2+10*2+11*2+6 {
		t.Errorf("keyByRune has %d entries, reference covers %d", len(keyByRune), 26*2+10*2+11*2+6)
	}
}

// TestFakerInputKeyDelegates verifies the FakerInput-only helpers derive from
// the unified KeyCodes tables.
func TestFakerInputKeyDelegates(t *testing.T) {
	code, mod, ok := FakerInputKeyFromRune('A')
	if !ok || code != fakerInputKeyA || mod != ModLShift {
		t.Errorf("FakerInputKeyFromRune('A') = (%#x, %#x, %v), want (%#x, %#x, true)", code, mod, ok, fakerInputKeyA, ModLShift)
	}

	code, ok = FakerInputKeyFromMakc(makc.KeyEnter)
	if !ok || code != fakerInputKeyEnter {
		t.Errorf("FakerInputKeyFromMakc(Enter) = (%#x, %v), want (%#x, true)", code, ok, fakerInputKeyEnter)
	}
}

// TestKeyCodesFromName verifies case-insensitive key-name resolution,
// including standalone modifier keys.
func TestKeyCodesFromName(t *testing.T) {
	cases := []struct {
		in   string
		want KeyCodes
		ok   bool
	}{
		{"enter", KeyCodes{makc.KeyEnter, fakerInputKeyEnter, 0}, true},
		{"ENTER", KeyCodes{makc.KeyEnter, fakerInputKeyEnter, 0}, true},
		{"f1", KeyCodes{makc.KeyF1, fakerInputKeyF1, 0}, true},
		{"F1", KeyCodes{makc.KeyF1, fakerInputKeyF1, 0}, true},
		{"left", KeyCodes{makc.KeyLeft, fakerInputKeyLeftArrow, 0}, true},
		{"ctrl", KeyCodes{makc.KeyControl, 0, ModLCtrl}, true},
		{"rctrl", KeyCodes{makc.KeyRightControl, 0, ModRCtrl}, true},
		{"shift", KeyCodes{makc.KeyShift, 0, ModLShift}, true},
		{"alt", KeyCodes{makc.KeyAlt, 0, ModLAlt}, true},
		{"gui", KeyCodes{makc.KeyLeftWindows, 0, ModLGui}, true},
		{"", KeyCodes{}, false},
		{"xyz", KeyCodes{}, false},
		{"a", KeyCodes{}, false}, // single runes are not names
	}
	for _, c := range cases {
		got, ok := KeyCodesFromName(c.in)
		if ok != c.ok || got != c.want {
			t.Errorf("KeyCodesFromName(%q) = (%+v, %v), want (%+v, %v)", c.in, got, ok, c.want, c.ok)
		}
	}
}

// TestKeyCodesFromRunes verifies rune-string resolution, including multi-rune
// chords and error cases.
func TestKeyCodesFromRunes(t *testing.T) {
	cases := []struct {
		in   string
		want []KeyCodes
		ok   bool
	}{
		{"a", []KeyCodes{{makc.KeyA, fakerInputKeyA, 0}}, true},
		{"A", []KeyCodes{{makc.KeyA, fakerInputKeyA, ModLShift}}, true},
		{"ab", []KeyCodes{
			{makc.KeyA, fakerInputKeyA, 0},
			{makc.KeyB, fakerInputKeyB, 0},
		}, true},
		{"AB", []KeyCodes{
			{makc.KeyA, fakerInputKeyA, ModLShift},
			{makc.KeyB, fakerInputKeyB, ModLShift},
		}, true},
		{"1!", []KeyCodes{
			{makc.Key1, fakerInputKey1, 0},
			{makc.Key1, fakerInputKey1, ModLShift},
		}, true},
		{"", nil, false},
		{"é", nil, false},
		{"aé", nil, false},
	}
	for _, c := range cases {
		got, ok := KeyCodesFromRunes(c.in)
		if ok != c.ok {
			t.Errorf("KeyCodesFromRunes(%q) ok = %v, want %v", c.in, ok, c.ok)
			continue
		}
		if !ok {
			continue
		}
		if len(got) != len(c.want) {
			t.Errorf("KeyCodesFromRunes(%q) len = %d, want %d", c.in, len(got), len(c.want))
			continue
		}
		for i := range got {
			if got[i] != c.want[i] {
				t.Errorf("KeyCodesFromRunes(%q)[%d] = %+v, want %+v", c.in, i, got[i], c.want[i])
			}
		}
	}
}
