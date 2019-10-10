package Node

import (
	"encoding/json"
	"fmt"
	"github.com/evscott/DistroA1/Models"
	"github.com/evscott/DistroA1/constants"
	"net"
	"time"
)

// Info contains all critical information that pertains to a nodes current state
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

// String is used to format a nodes key information; i.e. it's IP and Port
func (i Info) String() string {
	return "NodeInfo:{IP:" + i.IP + ", Port:" + i.Port + " }"
}

// Create is a constructor that is used to instantiate a new node with an IP, Port, and set of neighbors
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

// StartRST is the external command that triggers a node to create a distributed rooted spanning tree with it as the root
func (i *Info) StartRST() {
	i.Parent = i.Port
	i.Children = make(map[string]bool)
	i.ExpectedMsg = len(i.Neighbours)

	for _, neighbor := range i.Neighbours {
		msgOut := Models.Message{
			Source: i.Port,
			Intent: constants.IntentSendGo,
			Data:   "Some starting message",
		}

		if err := i.SendGo(msgOut, neighbor); err != nil {
			fmt.Println(err)
		}
	}
}

/**** Receive Handlers ****/

// ReceiveGo handles the event of a node receiving a "Go" messages on a distributed system
func (i *Info) ReceiveGo(msgIn Models.Message) {
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
			if err := i.SendBack(msgOut, i.Parent); err != nil {
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
					if err := i.SendGo(msgOut, neighbor); err != nil {
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
		if err := i.SendBack(msgOut, msgIn.Source); err != nil {
			fmt.Println(err)
		}
	}
}

// ReceiveBack handles the event of a node receiving a "Back" messages on a distributed system
func (i *Info) ReceiveBack(msgIn Models.Message) {
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
			if err := i.SendBack(msgOut, i.Parent); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("root: %v received ValSets received:", i.Port)
			for _, vs := range i.ValSet {
				fmt.Printf("%v\n", vs)
			}
		}
	}
}

/**** Send Handlers ****/

// SendPosition handles the event of a node sending a "Go" message to another node on a distributed system
func (i *Info) SendGo(msg Models.Message, dest string) error {
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

// SendPosition handles the event of a node sending a "Back" message to another node on a distributed system
func (i *Info) SendBack(msgIn Models.Message, dest string) error {
	connOut, err := net.DialTimeout("tcp", i.IP+":"+dest, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send back to %s:%s \n", i.IP, dest)
			return err
		}
	}

	if err := json.NewEncoder(connOut).Encode(&msgIn); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", msgIn)
		return err
	}
	return nil

}

/**** Node Communication Radar ****/

// ListenOnPort is intended to be a nodes satellite for receiving messages on a distributed system
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
			i.ReceiveGo(msg)
		case constants.IntentSendBack:
			i.ReceiveBack(msg)
		}
	}
}
