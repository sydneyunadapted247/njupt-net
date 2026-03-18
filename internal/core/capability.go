package core

// CapabilityLevel encodes certainty levels defined by SSOT.
type CapabilityLevel string

const (
	CapabilityConfirmed CapabilityLevel = "confirmed"
	CapabilityGuarded   CapabilityLevel = "guarded"
	CapabilityBlocked   CapabilityLevel = "blocked"
)
