package stock

import "github.com/ceheng.io/stock-go/internal/core"

// HostFallbackOptions configures fallback-host health tracking.
type HostFallbackOptions = core.HostFallbackOptions

// HostHealthStats exposes fallback-host health counters.
type HostHealthStats = core.HostHealthStats

// HostFallbackManager orders candidate hosts and tracks temporary cooldowns.
type HostFallbackManager = core.HostFallbackManager

// NewHostFallbackManager creates fallback-host governance.
func NewHostFallbackManager(options HostFallbackOptions) *HostFallbackManager {
	return core.NewHostFallbackManager(options)
}
