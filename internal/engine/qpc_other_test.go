//go:build !windows

package engine

import "time"

// On non-Windows platforms time.Now has nanosecond resolution, so use it
// directly. The tick frequency is expressed in nanoseconds per second to match
// the Windows QPC contract used by the latency benchmark.
func benchQPCTicksPerSecond() int64 { return int64(time.Second) }

func benchQPCNow() int64 { return time.Now().UnixNano() }
