package Node

import (
	"encoding/json"
	"fmt"
	"github.com/evscott/Distributed-RST/Models"
	"github.com/evscott/Distributed-RST/constants"
	"net"
	"time"
)

// Info represents the knowledge of a node participating in Raynal's definition of rooted spanning tree.
type Info struct {
	IP            string               `json:"IP"`
	Part          bool                 `json:"part"`
	Port          string               `json:"port"`
	Parent        string               `json:"parent"`
	Children      map[string]bool      `json:"children"`
	ExpectedMsg   int                  `json:"expectedMsg"`
	ProcKnown     map[string]bool      `json:"procKnown"`
	ChannelsKnown map[Models.Pair]bool `json:"channelsKnown"`
	Neighbours    []string             `json:"neighbors"`
	MessageInbox  map[string]bool      `json:"messageInbox"`
	ValSet        []Models.ValPair     `json:"valSet"`
}

// Just for pretty printing Node info
func (i Info) String() string {
	return "NodeInfo:{IP:" + i.IP + ", Port:" + i.Port + " }"
}

// Create is used a constructor that instantiates a new node using it's initial knowledge.
//
// A node must be created with initial knowledge of it's network IP, ID, and the IDs of it's neighbors.
func Create(ip, port string, neighbors []string) *Info {
	newNode := Info{
		IP:            ip,
		Port:          port,
		Neighbours:    neighbors,
		ChannelsKnown: make(map[Models.Pair]bool),
		ProcKnown:     make(map[string]bool),
		MessageInbox:  make(map[string]bool),
		ValSet: []Models.ValPair{{Node: port, Value: 1}},
	}

	newNode.ProcKnown[port] = true

	for _, neighbor := range neighbors {
		newNode.ChannelsKnown[Models.Pair{I: port, J: neighbor}] = true
	}

	return &newNode
}

// Start is an external command that triggers a node to contact its neighbors and begin the
// processes of creating a rooted spanning tree.
func (i *Info) Start() {
	i.Parent = i.Port
	i.Children = make(map[string]bool)
	i.ExpectedMsg = len(i.Neighbours)

	for _, neighbor := range i.Neighbours {
		msgOut := Models.Message{
			Source: i.Port,
			Intent: constants.IntentSendGo,
			Data:   "Some starting message",
		}

		if err := i.SendMsg(msgOut, neighbor); err != nil {
			fmt.Println(err)
		}
	}
}

// Go handles the event of a node receiving a "Go" messages with it's communication graph
// as part of creating a rooted spanning tree.
func (i *Info) Go(msgIn Models.Message) {
	if i.Parent == "" {
		i.Parent = msgIn.Source
		i.Children = make(map[string]bool)
		i.ExpectedMsg = len(i.Neighbours) - 1

		if i.ExpectedMsg == 0 {
			msgOut := Models.Message{
				Source: i.Port,
				Intent: constants.IntentSendBack,
				ValSet: i.ValSet,
			}
			if err := i.SendMsg(msgOut, i.Parent); err != nil { // send back
				fmt.Println(err)
			}
		} else {
			msgOut := Models.Message{
				Source: i.Port,
				Intent: constants.IntentSendGo,
				Data:   msgIn.Data,
			}
			for _, neighbor := range i.Neighbours {
				if neighbor != msgIn.Source {
					if err := i.SendMsg(msgOut, neighbor); err != nil { // Send go
						fmt.Println(err)
					}
				}
			}
		}
	} else {
		msgOut := Models.Message{
			Source: i.Port,
			Intent: constants.IntentSendBack,
			ValSet: nil,
		}
		if err := i.SendMsg(msgOut, msgIn.Source); err != nil { // Send back
			fmt.Println(err)
		}
	}
}

// Back handles the event of a node receiving a `Back` messages with it's communication graph
// as part of creating a rooted spanning tree.
func (i *Info) Back(msgIn Models.Message) {
	i.ExpectedMsg--
	if msgIn.ValSet != nil {
		i.Children[msgIn.Source] = true
		for _, valPair := range msgIn.ValSet {
			i.ValSet = append(i.ValSet, valPair)
		}
	}

	if i.ExpectedMsg == 0 { // Val_Set has been received by each child
		if i.Port != i.Parent {
			msgOut := Models.Message{
				Source: i.Port,
				Intent: constants.IntentSendBack,
				ValSet: i.ValSet,
			}
			if err := i.SendMsg(msgOut, i.Parent); err != nil { // Send back
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Root [%s] has received all ValSets\n", i.Port)
		}
	}
}

// SendPosition handles the event of a node sending a "Go" message to another node on a distributed system.
func (i *Info) SendMsg(msg Models.Message, dest string) error {
	connOut, err := net.DialTimeout("tcp", i.IP+":"+dest, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send go to %s:%s \n", i.IP, dest)
			return err
		}
	}

	if err := json.NewEncoder(connOut).Encode(&msg); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", msg)
		return err
	}
	return nil
}

// ListenOnPort is the communication satellite for a node that listens for incoming messages.
// Incoming messages are marshalled into a `Message` struct, and are directed to a handler
// depending on the messages `Intent`.
// Incoming messages that cannot be marshalled into a `Message` may cause erroneous behaviour.
func (i *Info) ListenOnPort() {
	ln, err := net.Listen("tcp", fmt.Sprint(":"+i.Port))
	if err != nil {
		fmt.Print(err)
		return
	}

	fmt.Printf("Staring node on %s:%s...\n", i.IP, i.Port)

	for {
		connIn, err := ln.Accept()
		if err != nil {
			if _, ok := err.(net.Error); ok {
				fmt.Printf("Error received while listening %s:%s \n", i.IP, i.Port)
			}
		}

		var msg Models.Message
		if err := json.NewDecoder(connIn).Decode(&msg); err != nil {
			fmt.Printf("Error decoding %v\n", err)
		}

		switch msg.Intent {
		case constants.IntentSendGo:
			i.Go(msg)
		case constants.IntentSendBack:
			i.Back(msg)
		}
	}
}
