//go:build windows

package engine

import (
	"syscall"
	"unsafe"
)

// Windows QueryPerformanceCounter is used for per-operation latency samples.
// time.Now is not good enough here: on some Windows hosts (including RDP/VM
// sessions) its granularity is ~0.5ms, so a microsecond-scale round trip is
// either measured as 0 or as a full clock tick. QPC typically ticks at 10MHz,
// giving ~100ns resolution.
var (
	procQueryPerformanceCounter   = kernel32.NewProc("QueryPerformanceCounter")
	procQueryPerformanceFrequency = kernel32.NewProc("QueryPerformanceFrequency")
)

// kernel32 is used for high-resolution latency sampling in benchmarks.
var kernel32 = syscall.NewLazyDLL("kernel32.dll")

func benchQPCTicksPerSecond() int64 {
	var freq int64
	_, _, _ = procQueryPerformanceFrequency.Call(uintptr(unsafe.Pointer(&freq)))
	return freq
}

func benchQPCNow() int64 {
	var ticks int64
	_, _, _ = procQueryPerformanceCounter.Call(uintptr(unsafe.Pointer(&ticks)))
	return ticks
}
