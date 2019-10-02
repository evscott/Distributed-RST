package Node

import (
	"encoding/json"
	"fmt"
	"github.com/evscott/DistroA1/Models"
	"github.com/evscott/DistroA1/constants"
	"net"
	"time"
)

/* TODO */
type Info struct {
	IP            string   `json:"IP"`
	Part          bool     `json:"part"`
	Port          string   `json:"port"`
	Parent        string   `json:"parent"`
	Children      map[string]bool   `json:"children"`
	ExpectedMsg   int      `json:"expectedMsg"`
	ProcKnown     map[string]bool  `json:"procKnown"`
	ChannelsKnown map[Models.Pair]bool   `json:"channelsKnown"`
	Neighbours    []string `json:"neighbors"`
	MessageInbox  map[string]bool `json:"messageInbox"`
	ValSet        []Models.ValPair   `json:"valSet"`
}

/* TODO */
func (i Info) String() string {
	return "NodeInfo:{IP:" + i.IP + ", Port:" + i.Port + " }"
}

/* TODO */
func Create(ip, port string, neighbors []string) *Info {
	newNode := Info{
		IP: ip,
		Port: port,
		Neighbours: neighbors,
		ChannelsKnown: make(map[Models.Pair]bool),
		ProcKnown: make(map[string]bool),
		MessageInbox: make(map[string]bool),
	}

	newNode.ProcKnown[port] = true

	for _, neighbor := range neighbors {
		newNode.ChannelsKnown[Models.Pair{I: port, J: neighbor}] = true
	}

	return &newNode
}

/* TODO */
func (i *Info) StartRST() {
	i.Parent = i.Port
	i.Children = make(map[string]bool)
	i.ExpectedMsg = len(i.Neighbours)

	for _, neighbor := range i.Neighbours {
		msgOut := Models.Message{
			Source: i.Port,
			Intent: constants.IntentSendGo,
			Data: "Some starting message",
		}

		if err := i.SendGo(msgOut, neighbor); err != nil {
			fmt.Println(err)
		}
	}
}

/* TODO */
func (i *Info) StartCG() {
	if i.Part == true {
		return
	}

	for _, n := range i.Neighbours {
		if err := i.SendPosition(i.Port, n, i.Neighbours); err != nil {
			fmt.Println("Error sending position", err)
			return
		}
	}

	i.Part = true
}

/**** Receive Handlers ****/

/* TODO */
func (i *Info) ReceivePosition(msg Models.Message) {
	if !i.Part {
		i.StartCG()
	}

	i.ProcKnown[msg.Source] = true

	count := 0
	for _, neighbor := range msg.Neighbours {
		if (i.ChannelsKnown[Models.Pair{I: msg.Source, J: neighbor}] || i.ChannelsKnown[Models.Pair{I: neighbor, J: msg.Source}]) {
			count++
		}
	}

	// If count is equal to the length of neighbors, this node already knows of the communication channels
	// between this messages source and neighbors -- i.e., this node has already received this message
	// and can therefore discard are as per 'The Forward/Discard Principle'
	//
	// If count is not equal to the length of neighbors, then all unknown communication channels should be added
	// to channels_known and the message should be forwarded to all of this nodes neighbors
	if count == len(msg.Neighbours) {
		return
	}

	// Add communication channels to list of channels known
	for _, neighbor := range msg.Neighbours {
		// If this communication channel is already known, ignore
		// Else, add to list of channels_known
		if (i.ChannelsKnown[Models.Pair{I: msg.Source, J: neighbor}] || i.ChannelsKnown[Models.Pair{I: neighbor, J: msg.Source}]) {
			continue
		}
		i.ChannelsKnown[Models.Pair{I: msg.Source, J: neighbor}] = true
	}

	// Forward received position to all other neighbours, excluding source of message
	for _, n := range i.Neighbours {
		if n == msg.Source {
			continue
		}
		if err := i.SendPosition(msg.Source, n, msg.Neighbours); err != nil {
			fmt.Println("Error sending position", err)
			return
		}
	}
}

/* TODO */
func (i *Info) ReceiveGo(msgIn Models.Message) {
	if i.Parent == "" {
		i.Parent = msgIn.Source
		i.Children = make(map[string]bool)
		i.ExpectedMsg = len(i.Neighbours)-1

		if i.ExpectedMsg == 0 {
			msgOut := Models.Message{
				Source: i.Port,
				Intent: constants.IntentSendBack,
				ValSet: []Models.ValPair{{Node: i.Port, Value: 1}},
			}
			if err := i.SendBack(msgOut); err != nil {
				fmt.Println(err)
			}
		} else {
			msgOut := Models.Message{
					Source: i.Port,
					Intent: constants.IntentSendGo,
					Data: msgIn.Data,
				}
			for _, neighbor := range i.Neighbours {
				if neighbor == msgIn.Source {
					continue
				}
				if err := i.SendGo(msgOut, neighbor); err != nil {
					fmt.Println(err)
				}
			}
		}
	}

	msgOut := Models.Message{
		Source: i.Port,
		Intent: constants.IntentSendBack,
		ValSet: i.ValSet,
		}
	if err := i.SendBack(msgOut); err != nil {
		fmt.Println(err)
	}
}

/* TODO */
func (i *Info) ReceiveBack(msgIn Models.Message) {
	i.ExpectedMsg--
	if len(i.ValSet) != 0 {
		i.Children[msgIn.Source] = true
	}
	for _, v := range msgIn.ValSet {
		i.ValSet = append(i.ValSet, v)
	}
	if i.ExpectedMsg == 0 {
		i.ValSet = append(i.ValSet, Models.ValPair{Node: i.Port, Value: 1})
		if i.Parent != i.Port {
			msgOut := Models.Message{
				Source:     i.Port,
				Intent:     constants.IntentSendBack,
			}
			if err := i.SendBack(msgOut); err != nil {
				fmt.Println(err)
			}
		} else {
			fmt.Printf("Root [%s] printing val sets: %v\n", i.Port, i.ValSet)
		}
	}
}

/**** Send Handlers ****/

/* TODO */
func (i *Info) SendPosition(source, dest string, neighbours []string) error {
	connOut, err := net.DialTimeout("tcp", i.IP+":"+dest, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send position to %s:%s \n", i.IP, dest)
			return err
		}
	}

	PositionMsg := Models.Message{Source: source, Intent: constants.IntentSendPosition, Neighbours: neighbours}

	if err := json.NewEncoder(connOut).Encode(&PositionMsg); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", PositionMsg)
		return err
	}
	return nil
}


/* TODO */
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

/* TODO */
func (i *Info) SendBack(msgIn Models.Message) error {
	connOut, err := net.DialTimeout("tcp", i.IP+":"+i.Parent, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send back to %s:%s \n", i.IP, i.Parent)
			return err
		}
	}

	if err := json.NewEncoder(connOut).Encode(&msgIn); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", msgIn)
		return err
	}
	return nil

}

/**** Sanity Check Handler *****/

/* TODO */
func (i *Info) SendPing(dest string) bool {
	connOut, err := net.DialTimeout("tcp", i.IP+":"+dest, time.Duration(10)*time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send ping to %s:%s \n", i.IP, dest)
			return false
		}
	}

	Ping := Models.Message{Source: i.Port, Intent: constants.IntentPing}
	if err := json.NewEncoder(connOut).Encode(&Ping); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", Ping)
		return false
	}
	return true
}

/**** Node Communication Radar ****/

/* TODO */
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
		case constants.IntentSendPosition:
			i.ReceivePosition(msg)
		case constants.IntentSendGo:
			i.ReceiveGo(msg)
		case constants.IntentSendBack:
			i.ReceiveBack(msg)
		case constants.IntentPing:
			fmt.Printf("Ping!\n")
		}
	}
}
