package Models

import (
	"github.com/evscott/DistroA1/constants"
)

/* TODO */
type Pair struct {
	I string `json:"i"`
	J string `json:"j"`
}

type ValPair struct {
	Node string `json:"i"`
	Value int `json:"j"`
}

/* A standard format for a Request/Response for adding node to cluster */
type Message struct {
	Source     string             `json:"source"`
	Intent     constants.Intent `json:"intent"`
	Neighbours []string         `json:"neighbours"`
	Data       string           	`json:"data"`
	ValSet     []ValPair    `json:"valSet"`
}

/* Just for pretty printing Request/Response info */
func (req Message) String() string {
	return "Message:{ Origin:" + req.Source + ", Intent: " + string(req.Intent) + " }\n"
}
