package vltemp

type VolumeState struct {
	Discs []string `json:"discs"`
}

type VolumeDiscoverNewDiscsUpdateResponse struct {
	// The workflow IDs of any newly-discovered Discs.
	Discovered []string `json:"discovered"`
}
