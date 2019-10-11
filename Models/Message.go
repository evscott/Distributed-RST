package Models

import (
	"github.com/evscott/Distributed-RST/constants"
)

// Pair representing a pair of node ID's, I and J.
type Pair struct {
	I string `json:"i"`
	J string `json:"j"`
}

// ValPair represents an element part of a ValSet, per Raynal's definition.
//
// Node represents the source of this ValPair.
// Value represents the result of a nodes computation from it's share of work in a rooted spanning tree collective task.
type ValPair struct {
	Node  string `json:"i"`
	Value int    `json:"j"`
}

// The format for a Request/Response in creating a rooted spanning tree.
//
// Source represents the Message sender.
// Intent represents the messages intent; i.e., whether it is to be handled by `Go` or `Back`.
// ValSet represents the ValSet being shared with a parent node, per Raynal's definition.
// Data represents some data to be shared using a message.
type Message struct {
	Source     string           `json:"source"`
	Intent     constants.Intent `json:"intent"`
	ValSet     []ValPair        `json:"valSet"`
	Data       string           `json:"data"`
}

// Just for pretty printing Request/Response info.
func (req Message) String() string {
	return "Message:{ Origin:" + req.Source + ", Intent: " + string(req.Intent) + " }\n"
}
