package input

import "github.com/aiwaki/makc"

// makcOpenPlan describes the makc.Open options for the primary backend choice
// and, when the selection is automatic, an optional fallback combination.
type makcOpenPlan struct {
	primary  []makc.Option
	fallback []makc.Option
}
