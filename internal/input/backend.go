package input

import "context"

// Backend is the platform-agnostic contract implemented by every input
// injection backend. The engine serializes calls per input class, so
// implementations do not need to be safe for concurrent use.
type Backend interface {
	// Name returns the backend name, for logging and diagnostics.
	Name() string
	// Press and Release inject one logical key press or release.
	Press(ctx context.Context, key KeyCodes) error
	Release(ctx context.Context, key KeyCodes) error
	// Mouse injects one backend-neutral mouse event.
	Mouse(ctx context.Context, ev MouseEvent) error
	// Close releases held input and backend resources. It must be safe to call
	// more than once.
	Close() error
}
