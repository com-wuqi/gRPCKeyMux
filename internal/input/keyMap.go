package input

import (
	"strings"

	"github.com/aiwaki/makc"
)

// Modifier key flags for FakerInputKeyboardReport.ShiftKeyFlags.
const (
	ModLCtrl  byte = 1
	ModLShift byte = 2
	ModLAlt   byte = 4
	ModLGui   byte = 8
	ModRCtrl  byte = 16
	ModRShift byte = 32
	ModRAlt   byte = 64
	ModRGui   byte = 128
)

// FakerInput USB HID keyboard usage codes (HUT 1.22).
const (
	fakerInputKeyNone byte = 0x00

	fakerInputKeyA byte = 0x04
	fakerInputKeyB byte = 0x05
	fakerInputKeyC byte = 0x06
	fakerInputKeyD byte = 0x07
	fakerInputKeyE byte = 0x08
	fakerInputKeyF byte = 0x09
	fakerInputKeyG byte = 0x0A
	fakerInputKeyH byte = 0x0B
	fakerInputKeyI byte = 0x0C
	fakerInputKeyJ byte = 0x0D
	fakerInputKeyK byte = 0x0E
	fakerInputKeyL byte = 0x0F
	fakerInputKeyM byte = 0x10
	fakerInputKeyN byte = 0x11
	fakerInputKeyO byte = 0x12
	fakerInputKeyP byte = 0x13
	fakerInputKeyQ byte = 0x14
	fakerInputKeyR byte = 0x15
	fakerInputKeyS byte = 0x16
	fakerInputKeyT byte = 0x17
	fakerInputKeyU byte = 0x18
	fakerInputKeyV byte = 0x19
	fakerInputKeyW byte = 0x1A
	fakerInputKeyX byte = 0x1B
	fakerInputKeyY byte = 0x1C
	fakerInputKeyZ byte = 0x1D

	fakerInputKey1 byte = 0x1E
	fakerInputKey2 byte = 0x1F
	fakerInputKey3 byte = 0x20
	fakerInputKey4 byte = 0x21
	fakerInputKey5 byte = 0x22
	fakerInputKey6 byte = 0x23
	fakerInputKey7 byte = 0x24
	fakerInputKey8 byte = 0x25
	fakerInputKey9 byte = 0x26
	fakerInputKey0 byte = 0x27

	fakerInputKeyEnter     byte = 0x28
	fakerInputKeyEscape    byte = 0x29
	fakerInputKeyBackspace byte = 0x2A
	fakerInputKeyTab       byte = 0x2B
	fakerInputKeySpace     byte = 0x2C
	fakerInputKeyMinus     byte = 0x2D
	fakerInputKeyEqual     byte = 0x2E
	fakerInputKeyLBracket  byte = 0x2F
	fakerInputKeyRBracket  byte = 0x30
	fakerInputKeyBackslash byte = 0x31
	fakerInputKeySemicolon byte = 0x33
	fakerInputKeyQuote     byte = 0x34
	fakerInputKeyBackquote byte = 0x35
	fakerInputKeyComma     byte = 0x36
	fakerInputKeyPeriod    byte = 0x37
	fakerInputKeySlash     byte = 0x38
	fakerInputKeyCapsLock  byte = 0x39

	fakerInputKeyF1  byte = 0x3A
	fakerInputKeyF2  byte = 0x3B
	fakerInputKeyF3  byte = 0x3C
	fakerInputKeyF4  byte = 0x3D
	fakerInputKeyF5  byte = 0x3E
	fakerInputKeyF6  byte = 0x3F
	fakerInputKeyF7  byte = 0x40
	fakerInputKeyF8  byte = 0x41
	fakerInputKeyF9  byte = 0x42
	fakerInputKeyF10 byte = 0x43
	fakerInputKeyF11 byte = 0x44
	fakerInputKeyF12 byte = 0x45

	fakerInputKeyPrintScreen byte = 0x46
	fakerInputKeyScrollLock  byte = 0x47
	fakerInputKeyInsert      byte = 0x49
	fakerInputKeyHome        byte = 0x4A
	fakerInputKeyPageUp      byte = 0x4B
	fakerInputKeyDelete      byte = 0x4C
	fakerInputKeyEnd         byte = 0x4D
	fakerInputKeyPageDown    byte = 0x4E
	fakerInputKeyRightArrow  byte = 0x4F
	fakerInputKeyLeftArrow   byte = 0x50
	fakerInputKeyDownArrow   byte = 0x51
	fakerInputKeyUpArrow     byte = 0x52
	fakerInputKeyNumLock     byte = 0x53
)

// KeyCodes holds one logical key in the two portable encodings used across the
// codebase: makc.Key (logical key; Windows virtual-key code, mapped to evdev by
// makc on Linux) and FakerInput (USB HID usage code). Modifiers carries the
// modifier flags (e.g. left shift) needed to produce the mapped rune.
//
// The Windows-only Interception scan code is derived from Makc by
// winputKey in keymap_windows.go, so it is not stored here.
type KeyCodes struct {
	Makc       makc.Key
	FakerInput byte
	Modifiers  byte
}

// fakerInputKeyTable is the single source of truth for logical keys (no
// modifiers), pairing the portable encodings. It is indexed by makc.Key through
// keyByMakc below; the Windows scan-code table is derived from the same keys.
var fakerInputKeyTable = []KeyCodes{
	// Letters
	{makc.KeyA, fakerInputKeyA, 0},
	{makc.KeyB, fakerInputKeyB, 0},
	{makc.KeyC, fakerInputKeyC, 0},
	{makc.KeyD, fakerInputKeyD, 0},
	{makc.KeyE, fakerInputKeyE, 0},
	{makc.KeyF, fakerInputKeyF, 0},
	{makc.KeyG, fakerInputKeyG, 0},
	{makc.KeyH, fakerInputKeyH, 0},
	{makc.KeyI, fakerInputKeyI, 0},
	{makc.KeyJ, fakerInputKeyJ, 0},
	{makc.KeyK, fakerInputKeyK, 0},
	{makc.KeyL, fakerInputKeyL, 0},
	{makc.KeyM, fakerInputKeyM, 0},
	{makc.KeyN, fakerInputKeyN, 0},
	{makc.KeyO, fakerInputKeyO, 0},
	{makc.KeyP, fakerInputKeyP, 0},
	{makc.KeyQ, fakerInputKeyQ, 0},
	{makc.KeyR, fakerInputKeyR, 0},
	{makc.KeyS, fakerInputKeyS, 0},
	{makc.KeyT, fakerInputKeyT, 0},
	{makc.KeyU, fakerInputKeyU, 0},
	{makc.KeyV, fakerInputKeyV, 0},
	{makc.KeyW, fakerInputKeyW, 0},
	{makc.KeyX, fakerInputKeyX, 0},
	{makc.KeyY, fakerInputKeyY, 0},
	{makc.KeyZ, fakerInputKeyZ, 0},

	// Digits
	{makc.Key0, fakerInputKey0, 0},
	{makc.Key1, fakerInputKey1, 0},
	{makc.Key2, fakerInputKey2, 0},
	{makc.Key3, fakerInputKey3, 0},
	{makc.Key4, fakerInputKey4, 0},
	{makc.Key5, fakerInputKey5, 0},
	{makc.Key6, fakerInputKey6, 0},
	{makc.Key7, fakerInputKey7, 0},
	{makc.Key8, fakerInputKey8, 0},
	{makc.Key9, fakerInputKey9, 0},

	// Whitespace and common control keys
	{makc.KeyEnter, fakerInputKeyEnter, 0},
	{makc.KeyEscape, fakerInputKeyEscape, 0},
	{makc.KeyBackspace, fakerInputKeyBackspace, 0},
	{makc.KeyTab, fakerInputKeyTab, 0},
	{makc.KeySpace, fakerInputKeySpace, 0},
	{makc.KeyCapsLock, fakerInputKeyCapsLock, 0},

	// Symbols
	{makc.KeyMinus, fakerInputKeyMinus, 0},
	{makc.KeyEquals, fakerInputKeyEqual, 0},
	{makc.KeyLeftSquareBracket, fakerInputKeyLBracket, 0},
	{makc.KeyRightSquareBracket, fakerInputKeyRBracket, 0},
	{makc.KeyBackslash, fakerInputKeyBackslash, 0},
	{makc.KeySemicolon, fakerInputKeySemicolon, 0},
	{makc.KeySingleQuote, fakerInputKeyQuote, 0},
	{makc.KeyBackQuote, fakerInputKeyBackquote, 0},
	{makc.KeyComma, fakerInputKeyComma, 0},
	{makc.KeyDot, fakerInputKeyPeriod, 0},
	{makc.KeySlash, fakerInputKeySlash, 0},

	// Function keys
	{makc.KeyF1, fakerInputKeyF1, 0},
	{makc.KeyF2, fakerInputKeyF2, 0},
	{makc.KeyF3, fakerInputKeyF3, 0},
	{makc.KeyF4, fakerInputKeyF4, 0},
	{makc.KeyF5, fakerInputKeyF5, 0},
	{makc.KeyF6, fakerInputKeyF6, 0},
	{makc.KeyF7, fakerInputKeyF7, 0},
	{makc.KeyF8, fakerInputKeyF8, 0},
	{makc.KeyF9, fakerInputKeyF9, 0},
	{makc.KeyF10, fakerInputKeyF10, 0},
	{makc.KeyF11, fakerInputKeyF11, 0},
	{makc.KeyF12, fakerInputKeyF12, 0},

	// Navigation and editing keys
	{makc.KeyScrollLock, fakerInputKeyScrollLock, 0},
	{makc.KeyInsert, fakerInputKeyInsert, 0},
	{makc.KeyHome, fakerInputKeyHome, 0},
	{makc.KeyPageUp, fakerInputKeyPageUp, 0},
	{makc.KeyDelete, fakerInputKeyDelete, 0},
	{makc.KeyEnd, fakerInputKeyEnd, 0},
	{makc.KeyPageDown, fakerInputKeyPageDown, 0},
	{makc.KeyRight, fakerInputKeyRightArrow, 0},
	{makc.KeyLeft, fakerInputKeyLeftArrow, 0},
	{makc.KeyDown, fakerInputKeyDownArrow, 0},
	{makc.KeyUp, fakerInputKeyUpArrow, 0},
	{makc.KeyNumLock, fakerInputKeyNumLock, 0},
}

// keyByMakc indexes fakerInputKeyTable by makc.Key (Windows virtual-key code).
var keyByMakc = func() map[makc.Key]KeyCodes {
	m := make(map[makc.Key]KeyCodes, len(fakerInputKeyTable))
	for _, kc := range fakerInputKeyTable {
		m[kc.Makc] = kc
	}
	return m
}()

// keyByRune maps a printable ASCII rune (or control rune) to its KeyCodes,
// including the modifier flags required for shifted variants.
var keyByRune = map[rune]KeyCodes{
	'a': {makc.KeyA, fakerInputKeyA, 0},
	'A': {makc.KeyA, fakerInputKeyA, ModLShift},
	'b': {makc.KeyB, fakerInputKeyB, 0},
	'B': {makc.KeyB, fakerInputKeyB, ModLShift},
	'c': {makc.KeyC, fakerInputKeyC, 0},
	'C': {makc.KeyC, fakerInputKeyC, ModLShift},
	'd': {makc.KeyD, fakerInputKeyD, 0},
	'D': {makc.KeyD, fakerInputKeyD, ModLShift},
	'e': {makc.KeyE, fakerInputKeyE, 0},
	'E': {makc.KeyE, fakerInputKeyE, ModLShift},
	'f': {makc.KeyF, fakerInputKeyF, 0},
	'F': {makc.KeyF, fakerInputKeyF, ModLShift},
	'g': {makc.KeyG, fakerInputKeyG, 0},
	'G': {makc.KeyG, fakerInputKeyG, ModLShift},
	'h': {makc.KeyH, fakerInputKeyH, 0},
	'H': {makc.KeyH, fakerInputKeyH, ModLShift},
	'i': {makc.KeyI, fakerInputKeyI, 0},
	'I': {makc.KeyI, fakerInputKeyI, ModLShift},
	'j': {makc.KeyJ, fakerInputKeyJ, 0},
	'J': {makc.KeyJ, fakerInputKeyJ, ModLShift},
	'k': {makc.KeyK, fakerInputKeyK, 0},
	'K': {makc.KeyK, fakerInputKeyK, ModLShift},
	'l': {makc.KeyL, fakerInputKeyL, 0},
	'L': {makc.KeyL, fakerInputKeyL, ModLShift},
	'm': {makc.KeyM, fakerInputKeyM, 0},
	'M': {makc.KeyM, fakerInputKeyM, ModLShift},
	'n': {makc.KeyN, fakerInputKeyN, 0},
	'N': {makc.KeyN, fakerInputKeyN, ModLShift},
	'o': {makc.KeyO, fakerInputKeyO, 0},
	'O': {makc.KeyO, fakerInputKeyO, ModLShift},
	'p': {makc.KeyP, fakerInputKeyP, 0},
	'P': {makc.KeyP, fakerInputKeyP, ModLShift},
	'q': {makc.KeyQ, fakerInputKeyQ, 0},
	'Q': {makc.KeyQ, fakerInputKeyQ, ModLShift},
	'r': {makc.KeyR, fakerInputKeyR, 0},
	'R': {makc.KeyR, fakerInputKeyR, ModLShift},
	's': {makc.KeyS, fakerInputKeyS, 0},
	'S': {makc.KeyS, fakerInputKeyS, ModLShift},
	't': {makc.KeyT, fakerInputKeyT, 0},
	'T': {makc.KeyT, fakerInputKeyT, ModLShift},
	'u': {makc.KeyU, fakerInputKeyU, 0},
	'U': {makc.KeyU, fakerInputKeyU, ModLShift},
	'v': {makc.KeyV, fakerInputKeyV, 0},
	'V': {makc.KeyV, fakerInputKeyV, ModLShift},
	'w': {makc.KeyW, fakerInputKeyW, 0},
	'W': {makc.KeyW, fakerInputKeyW, ModLShift},
	'x': {makc.KeyX, fakerInputKeyX, 0},
	'X': {makc.KeyX, fakerInputKeyX, ModLShift},
	'y': {makc.KeyY, fakerInputKeyY, 0},
	'Y': {makc.KeyY, fakerInputKeyY, ModLShift},
	'z': {makc.KeyZ, fakerInputKeyZ, 0},
	'Z': {makc.KeyZ, fakerInputKeyZ, ModLShift},

	'0': {makc.Key0, fakerInputKey0, 0},
	')': {makc.Key0, fakerInputKey0, ModLShift},
	'1': {makc.Key1, fakerInputKey1, 0},
	'!': {makc.Key1, fakerInputKey1, ModLShift},
	'2': {makc.Key2, fakerInputKey2, 0},
	'@': {makc.Key2, fakerInputKey2, ModLShift},
	'3': {makc.Key3, fakerInputKey3, 0},
	'#': {makc.Key3, fakerInputKey3, ModLShift},
	'4': {makc.Key4, fakerInputKey4, 0},
	'$': {makc.Key4, fakerInputKey4, ModLShift},
	'5': {makc.Key5, fakerInputKey5, 0},
	'%': {makc.Key5, fakerInputKey5, ModLShift},
	'6': {makc.Key6, fakerInputKey6, 0},
	'^': {makc.Key6, fakerInputKey6, ModLShift},
	'7': {makc.Key7, fakerInputKey7, 0},
	'&': {makc.Key7, fakerInputKey7, ModLShift},
	'8': {makc.Key8, fakerInputKey8, 0},
	'*': {makc.Key8, fakerInputKey8, ModLShift},
	'9': {makc.Key9, fakerInputKey9, 0},
	'(': {makc.Key9, fakerInputKey9, ModLShift},

	'`':  {makc.KeyBackQuote, fakerInputKeyBackquote, 0},
	'~':  {makc.KeyBackQuote, fakerInputKeyBackquote, ModLShift},
	'-':  {makc.KeyMinus, fakerInputKeyMinus, 0},
	'_':  {makc.KeyMinus, fakerInputKeyMinus, ModLShift},
	'=':  {makc.KeyEquals, fakerInputKeyEqual, 0},
	'+':  {makc.KeyEquals, fakerInputKeyEqual, ModLShift},
	'[':  {makc.KeyLeftSquareBracket, fakerInputKeyLBracket, 0},
	'{':  {makc.KeyLeftSquareBracket, fakerInputKeyLBracket, ModLShift},
	']':  {makc.KeyRightSquareBracket, fakerInputKeyRBracket, 0},
	'}':  {makc.KeyRightSquareBracket, fakerInputKeyRBracket, ModLShift},
	'\\': {makc.KeyBackslash, fakerInputKeyBackslash, 0},
	'|':  {makc.KeyBackslash, fakerInputKeyBackslash, ModLShift},
	';':  {makc.KeySemicolon, fakerInputKeySemicolon, 0},
	':':  {makc.KeySemicolon, fakerInputKeySemicolon, ModLShift},
	'\'': {makc.KeySingleQuote, fakerInputKeyQuote, 0},
	'"':  {makc.KeySingleQuote, fakerInputKeyQuote, ModLShift},
	',':  {makc.KeyComma, fakerInputKeyComma, 0},
	'<':  {makc.KeyComma, fakerInputKeyComma, ModLShift},
	'.':  {makc.KeyDot, fakerInputKeyPeriod, 0},
	'>':  {makc.KeyDot, fakerInputKeyPeriod, ModLShift},
	'/':  {makc.KeySlash, fakerInputKeySlash, 0},
	'?':  {makc.KeySlash, fakerInputKeySlash, ModLShift},

	' ':    {makc.KeySpace, fakerInputKeySpace, 0},
	'\n':   {makc.KeyEnter, fakerInputKeyEnter, 0},
	'\r':   {makc.KeyEnter, fakerInputKeyEnter, 0},
	'\t':   {makc.KeyTab, fakerInputKeyTab, 0},
	'\b':   {makc.KeyBackspace, fakerInputKeyBackspace, 0},
	'\x1b': {makc.KeyEscape, fakerInputKeyEscape, 0},
}

// KeyCodesFromRune returns the logical key (and required modifiers) for a
// printable ASCII rune or a supported control rune.
func KeyCodesFromRune(r rune) (KeyCodes, bool) {
	kc, ok := keyByRune[r]
	return kc, ok
}

// KeyCodesFromMakc returns the logical key for a makc.Key (Windows virtual-key
// code; makc maps it to the platform representation internally).
func KeyCodesFromMakc(k makc.Key) (KeyCodes, bool) {
	kc, ok := keyByMakc[k]
	return kc, ok
}

// FakerInputKeyFromRune maps a rune to its FakerInput USB HID key code and the
// modifier flags required to produce it.
func FakerInputKeyFromRune(r rune) (code byte, modifiers byte, ok bool) {
	kc, ok := keyByRune[r]
	return kc.FakerInput, kc.Modifiers, ok
}

// FakerInputKeyFromMakc converts a makc.Key (Windows virtual-key code) to a
// FakerInput USB HID key code.
func FakerInputKeyFromMakc(k makc.Key) (code byte, ok bool) {
	kc, ok := keyByMakc[k]
	return kc.FakerInput, ok
}

// keyByName maps canonical, lower-case key names (e.g. "enter", "f1", "left",
// "ctrl") to their KeyCodes. Modifier keys carry their modifier flag in
// Modifiers and a zero FakerInput code, because FakerInput encodes modifiers as
// report flags rather than key codes.
var keyByName = map[string]KeyCodes{
	// Whitespace and control keys
	"enter":     {makc.KeyEnter, fakerInputKeyEnter, 0},
	"return":    {makc.KeyEnter, fakerInputKeyEnter, 0},
	"escape":    {makc.KeyEscape, fakerInputKeyEscape, 0},
	"esc":       {makc.KeyEscape, fakerInputKeyEscape, 0},
	"backspace": {makc.KeyBackspace, fakerInputKeyBackspace, 0},
	"tab":       {makc.KeyTab, fakerInputKeyTab, 0},
	"space":     {makc.KeySpace, fakerInputKeySpace, 0},
	"capslock":  {makc.KeyCapsLock, fakerInputKeyCapsLock, 0},

	// Function keys
	"f1":  {makc.KeyF1, fakerInputKeyF1, 0},
	"f2":  {makc.KeyF2, fakerInputKeyF2, 0},
	"f3":  {makc.KeyF3, fakerInputKeyF3, 0},
	"f4":  {makc.KeyF4, fakerInputKeyF4, 0},
	"f5":  {makc.KeyF5, fakerInputKeyF5, 0},
	"f6":  {makc.KeyF6, fakerInputKeyF6, 0},
	"f7":  {makc.KeyF7, fakerInputKeyF7, 0},
	"f8":  {makc.KeyF8, fakerInputKeyF8, 0},
	"f9":  {makc.KeyF9, fakerInputKeyF9, 0},
	"f10": {makc.KeyF10, fakerInputKeyF10, 0},
	"f11": {makc.KeyF11, fakerInputKeyF11, 0},
	"f12": {makc.KeyF12, fakerInputKeyF12, 0},

	// Navigation and editing keys
	"scrolllock": {makc.KeyScrollLock, fakerInputKeyScrollLock, 0},
	"insert":     {makc.KeyInsert, fakerInputKeyInsert, 0},
	"ins":        {makc.KeyInsert, fakerInputKeyInsert, 0},
	"home":       {makc.KeyHome, fakerInputKeyHome, 0},
	"pageup":     {makc.KeyPageUp, fakerInputKeyPageUp, 0},
	"pgup":       {makc.KeyPageUp, fakerInputKeyPageUp, 0},
	"delete":     {makc.KeyDelete, fakerInputKeyDelete, 0},
	"del":        {makc.KeyDelete, fakerInputKeyDelete, 0},
	"end":        {makc.KeyEnd, fakerInputKeyEnd, 0},
	"pagedown":   {makc.KeyPageDown, fakerInputKeyPageDown, 0},
	"pgdn":       {makc.KeyPageDown, fakerInputKeyPageDown, 0},
	"right":      {makc.KeyRight, fakerInputKeyRightArrow, 0},
	"left":       {makc.KeyLeft, fakerInputKeyLeftArrow, 0},
	"down":       {makc.KeyDown, fakerInputKeyDownArrow, 0},
	"up":         {makc.KeyUp, fakerInputKeyUpArrow, 0},
	"numlock":    {makc.KeyNumLock, fakerInputKeyNumLock, 0},

	// Modifier keys. FakerInput has no key code for these; the modifier flag is
	// carried in Modifiers.
	"ctrl":       {makc.KeyControl, 0, ModLCtrl},
	"control":    {makc.KeyControl, 0, ModLCtrl},
	"lctrl":      {makc.KeyLeftControl, 0, ModLCtrl},
	"leftctrl":   {makc.KeyLeftControl, 0, ModLCtrl},
	"rctrl":      {makc.KeyRightControl, 0, ModRCtrl},
	"rightctrl":  {makc.KeyRightControl, 0, ModRCtrl},
	"shift":      {makc.KeyShift, 0, ModLShift},
	"lshift":     {makc.KeyLeftShift, 0, ModLShift},
	"leftshift":  {makc.KeyLeftShift, 0, ModLShift},
	"rshift":     {makc.KeyRightShift, 0, ModRShift},
	"rightshift": {makc.KeyRightShift, 0, ModRShift},
	"alt":        {makc.KeyAlt, 0, ModLAlt},
	"lalt":       {makc.KeyLeftAlt, 0, ModLAlt},
	"leftalt":    {makc.KeyLeftAlt, 0, ModLAlt},
	"ralt":       {makc.KeyRightAlt, 0, ModRAlt},
	"rightalt":   {makc.KeyRightAlt, 0, ModRAlt},
	"gui":        {makc.KeyLeftWindows, 0, ModLGui},
	"win":        {makc.KeyLeftWindows, 0, ModLGui},
	"windows":    {makc.KeyLeftWindows, 0, ModLGui},
	"super":      {makc.KeyLeftWindows, 0, ModLGui},
	"meta":       {makc.KeyLeftWindows, 0, ModLGui},
}

// KeyCodesFromName resolves a key name (e.g. "enter", "F1", "left", "ctrl")
// case-insensitively via keyByName.
func KeyCodesFromName(name string) (KeyCodes, bool) {
	if name == "" {
		return KeyCodes{}, false
	}
	kc, ok := keyByName[strings.ToLower(name)]
	return kc, ok
}

// KeyCodesFromRunes resolves a rune string into an ordered slice of KeyCodes,
// one per rune. Every rune must be mappable; an empty string or any unknown
// rune yields ok=false. Multiple runes represent a simultaneous chord.
func KeyCodesFromRunes(s string) ([]KeyCodes, bool) {
	if s == "" {
		return nil, false
	}
	runes := []rune(s)
	keys := make([]KeyCodes, 0, len(runes))
	for _, r := range runes {
		kc, ok := KeyCodesFromRune(r)
		if !ok {
			return nil, false
		}
		keys = append(keys, kc)
	}
	return keys, true
}

// SamePhysicalKey reports whether two KeyCodes refer to the same physical key,
// ignoring modifier flags. The FakerInput code is used as the backend-neutral
// identity because every non-modifier key maps to it one-to-one; modifier-only
// keys (code 0) are identified by their modifier flag instead.
func SamePhysicalKey(a, b KeyCodes) bool {
	if a.FakerInput != 0 || b.FakerInput != 0 {
		return a.FakerInput == b.FakerInput
	}
	return a.Modifiers == b.Modifiers
}
