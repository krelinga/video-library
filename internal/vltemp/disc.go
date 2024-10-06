package vltemp

const (
	Disc                = "disc"
	DiscUpdateBootstrap = "disc-update-bootstrap"
)


type DiscState struct {
	Videos []string `json:"videos"`
}
