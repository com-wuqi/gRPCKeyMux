//go:build windows

package input

import (
	"github.com/aiwaki/makc"
	"github.com/rpdg/winput"
)

// scanByMakc maps the portable logical key (makc.Key) to the Interception
// hardware scan code used by the winput backends. It is kept on the Windows
// build only so that the shared key table and the rest of the input package
// stay portable.
var scanByMakc = map[makc.Key]winput.Key{
	// Letters
	makc.KeyA: winput.KeyA,
	makc.KeyB: winput.KeyB,
	makc.KeyC: winput.KeyC,
	makc.KeyD: winput.KeyD,
	makc.KeyE: winput.KeyE,
	makc.KeyF: winput.KeyF,
	makc.KeyG: winput.KeyG,
	makc.KeyH: winput.KeyH,
	makc.KeyI: winput.KeyI,
	makc.KeyJ: winput.KeyJ,
	makc.KeyK: winput.KeyK,
	makc.KeyL: winput.KeyL,
	makc.KeyM: winput.KeyM,
	makc.KeyN: winput.KeyN,
	makc.KeyO: winput.KeyO,
	makc.KeyP: winput.KeyP,
	makc.KeyQ: winput.KeyQ,
	makc.KeyR: winput.KeyR,
	makc.KeyS: winput.KeyS,
	makc.KeyT: winput.KeyT,
	makc.KeyU: winput.KeyU,
	makc.KeyV: winput.KeyV,
	makc.KeyW: winput.KeyW,
	makc.KeyX: winput.KeyX,
	makc.KeyY: winput.KeyY,
	makc.KeyZ: winput.KeyZ,

	// Digits
	makc.Key0: winput.Key0,
	makc.Key1: winput.Key1,
	makc.Key2: winput.Key2,
	makc.Key3: winput.Key3,
	makc.Key4: winput.Key4,
	makc.Key5: winput.Key5,
	makc.Key6: winput.Key6,
	makc.Key7: winput.Key7,
	makc.Key8: winput.Key8,
	makc.Key9: winput.Key9,

	// Whitespace and common control keys
	makc.KeyEnter:     winput.KeyEnter,
	makc.KeyEscape:    winput.KeyEsc,
	makc.KeyBackspace: winput.KeyBkSp,
	makc.KeyTab:       winput.KeyTab,
	makc.KeySpace:     winput.KeySpace,
	makc.KeyCapsLock:  winput.KeyCaps,

	// Symbols
	makc.KeyMinus:              winput.KeyMinus,
	makc.KeyEquals:             winput.KeyEqual,
	makc.KeyLeftSquareBracket:  winput.KeyLBr,
	makc.KeyRightSquareBracket: winput.KeyRBr,
	makc.KeyBackslash:          winput.KeyBackslash,
	makc.KeySemicolon:          winput.KeySemi,
	makc.KeySingleQuote:        winput.KeyQuot,
	makc.KeyBackQuote:          winput.KeyTick,
	makc.KeyComma:              winput.KeyComma,
	makc.KeyDot:                winput.KeyDot,
	makc.KeySlash:              winput.KeySlash,

	// Function keys
	makc.KeyF1:  winput.KeyF1,
	makc.KeyF2:  winput.KeyF2,
	makc.KeyF3:  winput.KeyF3,
	makc.KeyF4:  winput.KeyF4,
	makc.KeyF5:  winput.KeyF5,
	makc.KeyF6:  winput.KeyF6,
	makc.KeyF7:  winput.KeyF7,
	makc.KeyF8:  winput.KeyF8,
	makc.KeyF9:  winput.KeyF9,
	makc.KeyF10: winput.KeyF10,
	makc.KeyF11: winput.KeyF11,
	makc.KeyF12: winput.KeyF12,

	// Navigation and editing keys
	makc.KeyScrollLock: winput.KeyScroll,
	makc.KeyInsert:     winput.KeyInsert,
	makc.KeyHome:       winput.KeyHome,
	makc.KeyPageUp:     winput.KeyPageUp,
	makc.KeyDelete:     winput.KeyDelete,
	makc.KeyEnd:        winput.KeyEnd,
	makc.KeyPageDown:   winput.KeyPageDown,
	makc.KeyRight:      winput.KeyRight,
	makc.KeyLeft:       winput.KeyLeft,
	makc.KeyDown:       winput.KeyArrowDown,
	makc.KeyUp:         winput.KeyArrowUp,
	makc.KeyNumLock:    winput.KeyNumLock,

	// Modifiers. winput cannot distinguish left/right variants for these, and
	// the Windows/Gui key is not representable at all (lookup returns 0).
	makc.KeyControl:      winput.KeyCtrl,
	makc.KeyLeftControl:  winput.KeyCtrl,
	makc.KeyRightControl: winput.KeyCtrl,
	makc.KeyShift:        winput.KeyShift,
	makc.KeyLeftShift:    winput.KeyShift,
	makc.KeyRightShift:   winput.KeyShift,
	makc.KeyAlt:          winput.KeyAlt,
	makc.KeyLeftAlt:      winput.KeyAlt,
	makc.KeyRightAlt:     winput.KeyAlt,
}

// keyByWinput indexes the logical key table by Interception scan code. Only keys
// present in fakerInputKeyTable are included, matching the historical behavior
// (standalone modifiers are not resolvable by scan code).
var keyByWinput = func() map[winput.Key]KeyCodes {
	m := make(map[winput.Key]KeyCodes, len(fakerInputKeyTable))
	for _, kc := range fakerInputKeyTable {
		if scan, ok := scanByMakc[kc.Makc]; ok {
			m[scan] = kc
		}
	}
	return m
}()

// winputKey returns the Interception scan code for a logical key. A zero result
// means the key cannot be represented by the winput backends.
func winputKey(k KeyCodes) winput.Key {
	return scanByMakc[k.Makc]
}

// KeyCodesFromWinput returns the logical key for a winput.Key (Interception
// scan code).
func KeyCodesFromWinput(k winput.Key) (KeyCodes, bool) {
	kc, ok := keyByWinput[k]
	return kc, ok
}

// FakerInputKeyFromWinput converts a winput.Key (Interception scan code) to a
// FakerInput USB HID key code.
func FakerInputKeyFromWinput(k winput.Key) (code byte, ok bool) {
	kc, ok := keyByWinput[k]
	return kc.FakerInput, ok
}
