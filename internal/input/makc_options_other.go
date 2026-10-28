//go:build !windows && !linux

package input

// makcOpenPlanForPlatform uses makc's defaults on platforms other than Windows
// and Linux. Linux-specific backend selection lives in makc_options_linux.go.
func makcOpenPlanForPlatform() makcOpenPlan {
	return makcOpenPlan{}
}
