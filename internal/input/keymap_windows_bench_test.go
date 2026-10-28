//go:build windows

package input

import (
	"testing"

	"github.com/rpdg/winput"
)

func BenchmarkKeyCodesFromWinput(b *testing.B) {
	a, ok := KeyCodesFromName("space")
	if !ok {
		b.Fatal("key space is missing from the table")
	}
	scan := winputKey(a)
	if scan == 0 {
		b.Fatal("key space has no scan code")
	}
	b.Run("hit", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			kc, ok := KeyCodesFromWinput(scan)
			if !ok {
				b.Fatal("winput lookup miss")
			}
			benchSinkKey = kc
		}
	})
	b.Run("miss", func(b *testing.B) {
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			kc, ok := KeyCodesFromWinput(winput.Key(0xFFFF))
			if ok {
				b.Fatal("unexpected winput lookup hit")
			}
			benchSinkKey = kc
		}
	})
}
