//go:build windows

package input

// makcOpenPlanForPlatform keeps makc's Windows defaults untouched: no options
// are passed, so the platform's automatic backend selection is preserved.
func makcOpenPlanForPlatform() makcOpenPlan {
	return makcOpenPlan{}
}
