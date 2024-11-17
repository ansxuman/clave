package objects

type InitResult struct {
	NeedsOnboarding   bool `json:"needsOnboarding"`
	NeedsVerification bool `json:"needsVerification"`
}
