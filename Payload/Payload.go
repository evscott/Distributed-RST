package Payload

import (
	"github.com/evscott/DistroA1/constants"
)

/* A standard format for a Request/Response for adding node to cluster */
type Payload struct {
	Source   string             `json:"source"`
	Dest     string             `json:"dest"`
	Intent   constants.Intent   `json:"intent"`
	Message []string            `json:"message"`
}

/* Just for pretty printing Request/Response info */
func (req Payload) String() string {
	return "Message:{ Source:" + req.Source + ", Intent: " + string(req.Intent) + " }\n"
}
