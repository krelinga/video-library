package vltypes

type VideoLineage struct {
	FromDisc *VideoFromDisc `json:"from_disc"`
	// TODO: eventually support other options here.
}

type VideoFromDisc struct {
	DiscID   string `json:"disc_id"`
	Filename string `json:"filename"`
}

type VideoUpdateBootstrapRequest struct {
	Lineage *VideoLineage `json:"lineage"`
}
