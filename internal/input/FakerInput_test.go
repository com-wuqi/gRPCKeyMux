//go:build windows

package input

import (
	"bytes"
	"testing"
)

func TestPutUint32LE(t *testing.T) {
	b := make([]byte, 4)
	putUint32LE(b, 0x01020304)
	if !bytes.Equal(b, []byte{0x04, 0x03, 0x02, 0x01}) {
		t.Fatalf("unexpected little-endian encoding: %v", b)
	}
}

func TestFillKeyboardReport(t *testing.T) {
	var got [fakerInputKeyboardReportSize]byte
	fillKeyboardReport(&got, ModLCtrl|ModLShift, [fakerInputKeyCodeCount]byte{fakerInputKeyA, fakerInputKeyB, 0, 0, 0, 0})
	want := [fakerInputKeyboardReportSize]byte{0x01, ModLCtrl | ModLShift, 0x00, fakerInputKeyA, fakerInputKeyB, 0, 0, 0, 0}
	if got != want {
		t.Fatalf("fillKeyboardReport = %v, want %v", got, want)
	}
}

func TestFillRelativeMouseReport(t *testing.T) {
	var got [fakerInputRelativeMouseReportSize]byte
	fillRelativeMouseReport(&got, fakerInputMouseButtonLeft|fakerInputMouseButtonMiddle, -2, 3, -4, 5)
	want := [fakerInputRelativeMouseReportSize]byte{
		0x03,
		fakerInputMouseButtonLeft | fakerInputMouseButtonMiddle,
		0xFE, 0xFF, // XValue = -2
		0x03, 0x00, // YValue = 3
		0xFC, // WheelPosition = -4
		0x05, // HWheelPosition = 5
	}
	if got != want {
		t.Fatalf("fillRelativeMouseReport = %v, want %v", got, want)
	}
}

func TestFillAbsoluteMouseReport(t *testing.T) {
	var got [fakerInputAbsoluteMouseReportSize]byte
	fillAbsoluteMouseReport(&got, fakerInputMouseButtonRight, 0x1234, 0xABCD, -7)
	want := [fakerInputAbsoluteMouseReportSize]byte{
		0x04,
		fakerInputMouseButtonRight,
		0x34, 0x12, // XValue = 0x1234
		0xCD, 0xAB, // YValue = 0xABCD
		0xF9, // WheelPosition = -7
	}
	if got != want {
		t.Fatalf("fillAbsoluteMouseReport = %v, want %v", got, want)
	}
}

// TestMouseButtonState verifies mouse buttons accumulate as a bitmask and
// release independently.
func TestMouseButtonState(t *testing.T) {
	d := &FakerInputDevice{}
	d.MouseButtonDown(fakerInputMouseButtonLeft)
	d.MouseButtonDown(fakerInputMouseButtonRight)

	if d.mouseButtons != fakerInputMouseButtonLeft|fakerInputMouseButtonRight {
		t.Fatalf("mouseButtons = %#x, want %#x", d.mouseButtons, fakerInputMouseButtonLeft|fakerInputMouseButtonRight)
	}

	d.MouseButtonUp(fakerInputMouseButtonLeft)
	if d.mouseButtons != fakerInputMouseButtonRight {
		t.Fatalf("mouseButtons = %#x, want %#x", d.mouseButtons, fakerInputMouseButtonRight)
	}

	d.ReleaseAllButtons()
	if d.mouseButtons != 0 {
		t.Fatalf("mouseButtons = %#x, want 0", d.mouseButtons)
	}
}

// TestMouseMoveCarriesHeldButtons verifies that a movement report built from
// the device's state still carries the currently held buttons (required for
// dragging).
func TestMouseMoveCarriesHeldButtons(t *testing.T) {
	d := &FakerInputDevice{}
	d.MouseButtonDown(fakerInputMouseButtonLeft)

	var got [fakerInputRelativeMouseReportSize]byte
	fillRelativeMouseReport(&got, d.mouseButtons, 10, 20, 0, 0)
	if got[1] != fakerInputMouseButtonLeft {
		t.Fatalf("move report button byte = %#x, want %#x", got[1], fakerInputMouseButtonLeft)
	}
}

// TestMouseKeyboardIndependent verifies keyboard and mouse state are fully
// independent: holding a key and a mouse button simultaneously does not disturb
// either state. This is the core of the simultaneous input guarantee.
func TestMouseKeyboardIndependent(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, ModLShift)
	d.MouseButtonDown(fakerInputMouseButtonLeft)

	mods, keys := d.heldState()
	if mods != ModLShift {
		t.Fatalf("keyboard mods = %#x, want %#x", mods, ModLShift)
	}
	if !containsKey(keys, fakerInputKeyA) {
		t.Fatalf("held keys = %v, want a", keys)
	}

	if d.mouseButtons != fakerInputMouseButtonLeft {
		t.Fatalf("mouseButtons = %#x, want %#x", d.mouseButtons, fakerInputMouseButtonLeft)
	}

	d.MouseButtonUp(fakerInputMouseButtonLeft)
	d.keyUp(fakerInputKeyA, ModLShift)

	if d.mouseButtons != 0 {
		t.Fatalf("mouseButtons = %#x, want 0", d.mouseButtons)
	}
	if d.keyCount != 0 {
		t.Fatalf("keyCount = %d, want 0", d.keyCount)
	}
}

func TestFillControlReportLayout(t *testing.T) {
	inner := []byte{0x01, 0x00, 0x00, 0x04, 0x00, 0x00, 0x00, 0x00, 0x00}
	var report [65]byte
	fillControlReport(&report, inner)
	if report[0] != fakerInputReportIDControl {
		t.Fatalf("unexpected report ID: %#x, want %#x", report[0], fakerInputReportIDControl)
	}
	if report[1] != byte(len(inner)) {
		t.Fatalf("unexpected ReportLength: %d, want %d", report[1], len(inner))
	}
	if !bytes.Equal(report[2:2+len(inner)], inner) {
		t.Fatalf("inner report mismatch: %v", report[2:])
	}
	if !bytes.Equal(report[2+len(inner):], make([]byte, 65-2-len(inner))) {
		t.Fatalf("padding is not zeroed: %v", report[2+len(inner):])
	}
}

// TestKeyDownAccumulates verifies that pressing multiple distinct keys holds
// them all simultaneously (press a, press b = hold ab).
func TestKeyDownAccumulates(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, 0)
	d.keyDown(fakerInputKeyB, 0)

	if d.keyCount != 2 {
		t.Fatalf("keyCount = %d, want 2", d.keyCount)
	}
	mods, keys := d.heldState()
	if mods != 0 {
		t.Fatalf("mods = %#x, want 0", mods)
	}
	if !containsKey(keys, fakerInputKeyA) || !containsKey(keys, fakerInputKeyB) {
		t.Fatalf("held keys = %v, want a and b", keys)
	}
}

// TestKeyDownIdempotent verifies pressing the same key twice does not create a
// duplicate slot or double up.
func TestKeyDownIdempotent(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, 0)
	d.keyDown(fakerInputKeyA, 0)

	if d.keyCount != 1 {
		t.Fatalf("keyCount = %d, want 1", d.keyCount)
	}
}

// TestKeyUpReleasesSingle verifies releasing one key leaves the other held.
func TestKeyUpReleasesSingle(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, 0)
	d.keyDown(fakerInputKeyB, 0)
	d.keyUp(fakerInputKeyA, 0)

	if d.keyCount != 1 {
		t.Fatalf("keyCount = %d, want 1", d.keyCount)
	}
	_, keys := d.heldState()
	if containsKey(keys, fakerInputKeyA) || !containsKey(keys, fakerInputKeyB) {
		t.Fatalf("held keys = %v, want only b", keys)
	}
}

// TestModifierUnion verifies modifiers are reference-counted: releasing a
// shifted key does not drop shift while another shifted key is still held.
func TestModifierUnion(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, ModLShift)
	d.keyDown(fakerInputKey2, ModLShift) // '@' shares shift
	d.keyUp(fakerInputKeyA, ModLShift)

	mods, keys := d.heldState()
	if mods&ModLShift == 0 {
		t.Fatalf("shift was dropped while another shifted key is still held: mods=%#x", mods)
	}
	if containsKey(keys, fakerInputKeyA) || !containsKey(keys, fakerInputKey2) {
		t.Fatalf("held keys = %v, want only 2", keys)
	}
}

// TestKeyUpClearsModifiersOnEmptyMask verifies a key stays held while it still
// has a modifier mask, and is removed once its mask is empty.
func TestKeyUpClearsModifiersOnEmptyMask(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, ModLShift)
	d.keyDown(fakerInputKeyA, ModLCtrl) // accumulates ctrl on the same key
	if d.keyCount != 1 {
		t.Fatalf("keyCount = %d, want 1", d.keyCount)
	}

	d.keyUp(fakerInputKeyA, ModLCtrl)
	if d.keyCount != 1 {
		t.Fatalf("key still should be held: keyCount = %d, want 1", d.keyCount)
	}
	mods, _ := d.heldState()
	if mods != ModLShift {
		t.Fatalf("mods = %#x, want %#x", mods, ModLShift)
	}

	d.keyUp(fakerInputKeyA, ModLShift)
	if d.keyCount != 0 {
		t.Fatalf("key should be released: keyCount = %d, want 0", d.keyCount)
	}
}

// TestKeyDownSixKeyCap verifies the held set is capped at 6 keys.
func TestKeyDownSixKeyCap(t *testing.T) {
	d := &FakerInputDevice{}
	codes := []byte{
		fakerInputKeyA, fakerInputKeyB, fakerInputKeyC, fakerInputKeyD,
		fakerInputKeyE, fakerInputKeyF, fakerInputKeyG,
	}
	for _, c := range codes {
		d.keyDown(c, 0)
	}
	if d.keyCount != fakerInputKeyCodeCount {
		t.Fatalf("keyCount = %d, want %d", d.keyCount, fakerInputKeyCodeCount)
	}
}

// TestReleaseAll verifies ReleaseAll clears the held-key state.
func TestReleaseAll(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(fakerInputKeyA, ModLShift)
	d.keyDown(fakerInputKeyB, ModLCtrl)
	d.releaseAll()

	if d.keyCount != 0 {
		t.Fatalf("keyCount = %d, want 0", d.keyCount)
	}
	mods, keys := d.heldState()
	if mods != 0 {
		t.Fatalf("mods = %#x, want 0", mods)
	}
	if keys != [fakerInputKeyCodeCount]byte{} {
		t.Fatalf("keys = %v, want empty", keys)
	}
}

func containsKey(keys [fakerInputKeyCodeCount]byte, code byte) bool {
	for _, k := range keys {
		if k == code {
			return true
		}
	}
	return false
}

// TestStandaloneModifier verifies that modifier-only presses/releases are
// tracked independently of held keys and survive key press/release around them.
func TestStandaloneModifier(t *testing.T) {
	d := &FakerInputDevice{}

	d.keyDown(0, ModLCtrl)
	mods, _ := d.heldState()
	if mods&ModLCtrl == 0 {
		t.Fatalf("ctrl not held: mods=%#x", mods)
	}

	d.keyDown(fakerInputKeyA, 0)
	mods, keys := d.heldState()
	if mods&ModLCtrl == 0 {
		t.Fatalf("ctrl lost after pressing a key: mods=%#x", mods)
	}
	if !containsKey(keys, fakerInputKeyA) {
		t.Fatalf("key a not held: %v", keys)
	}

	d.keyUp(fakerInputKeyA, 0)
	mods, _ = d.heldState()
	if mods&ModLCtrl == 0 {
		t.Fatalf("ctrl lost after releasing a key: mods=%#x", mods)
	}

	d.keyUp(0, ModLCtrl)
	mods, _ = d.heldState()
	if mods != 0 {
		t.Fatalf("mods = %#x, want 0", mods)
	}
}

// TestReleaseAllClearsStandaloneModifiers verifies ReleaseAll clears held
// modifier-only state as well as held keys.
func TestReleaseAllClearsStandaloneModifiers(t *testing.T) {
	d := &FakerInputDevice{}
	d.keyDown(0, ModLAlt)
	d.keyDown(fakerInputKeyA, 0)
	d.releaseAll()

	if d.keyCount != 0 {
		t.Fatalf("keyCount = %d, want 0", d.keyCount)
	}
	mods, _ := d.heldState()
	if mods != 0 {
		t.Fatalf("mods = %#x, want 0", mods)
	}
}

func TestFakerInputAbsCoordinate(t *testing.T) {
	cases := []struct {
		in   int32
		want uint16
	}{
		{0, 0},
		{10000, 32767},
		{5000, 16383},
		{-1, 0},
		{10001, 32767},
	}
	for _, c := range cases {
		if got := FakerInputAbsCoordinate(c.in); got != c.want {
			t.Errorf("FakerInputAbsCoordinate(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}

func TestFakerInputRelDelta(t *testing.T) {
	cases := []struct {
		in   int32
		want int16
	}{
		{0, 0},
		{1, 1},
		{-1, -1},
		{32767, 32767},
		{-32767, -32767},
		{32768, 32767},
		{-32768, -32767},
		{40000, 32767},
		{-40000, -32767},
	}
	for _, c := range cases {
		if got := FakerInputRelDelta(c.in); got != c.want {
			t.Errorf("FakerInputRelDelta(%d) = %d, want %d", c.in, got, c.want)
		}
	}
}
