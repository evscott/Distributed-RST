package Models

import (
	"github.com/evscott/Distributed-RST/constants"
)

// Pair is a generic struct representing a pair of two Node ID's, I and J
type Pair struct {
	I string `json:"i"`
	J string `json:"j"`
}

// ValPair is a Node ID and Value pairing used for ValSet communication
type ValPair struct {
	Node  string `json:"i"`
	Value int    `json:"j"`
}

// A standard format for a Request/Response for adding node to cluster
type Message struct {
	Source     string           `json:"source"`
	Intent     constants.Intent `json:"intent"`
	Data       string           `json:"data"`
	ValSet     []ValPair        `json:"valSet"`
}

// Just for pretty printing Request/Response info
func (req Message) String() string {
	return "Message:{ Origin:" + req.Source + ", Intent: " + string(req.Intent) + " }\n"
}
