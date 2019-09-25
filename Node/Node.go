package Node

import (
	"encoding/json"
	"fmt"
	"github.com/evscott/DistroA1/Payload"
	"github.com/evscott/DistroA1/constants"
	"net"
	"time"
)

type Pair struct {
	I string  `json:"i"`
	J string  `json:"j"`
}

/* Information/Metadata about Node */
type Info struct {
	IP 		      string 			`json:"IP"`
	Part          bool              `json:"Part"`
	Port          string 			`json:"Port"`
	Parent        string 			`json:"Parent"`
	Children   	  []string 			`json:"Children"`
	ExpectedMsg   int			    `json:"ExpectedMsg"`
	ProcKnown     []Pair			`json:"ProcKnown"`
	ChannelsKnown []Pair		    `json:"ChannelsKnown"`
	Neighbours    []string			`json:"Neighbors"`
}

func (i Info) String() string {
	return "NodeInfo:{IP:" + i.IP + ", Port:" + i.Port + " }"
}

func Create(ip, port string, neighbors []string) *Info {
	newNode := Info{IP: ip, Port: port, Neighbours: neighbors}
	return &newNode
}

func (i *Info) Start() {
	if i.Part == true {
		return
	}

	for _, n := range i.Neighbours {
		if err := i.SendPosition(n); err != nil {
			fmt.Println("Error sending position", err)
			return
		}
	}

	i.Part = true
}

/**** Receive Handlers ****/

func (i *Info) ReceivePosition(source string, neighbors []string) {
	for _, n := range neighbors {
		i.ChannelsKnown = append(i.ChannelsKnown, Pair{I: source, J: n})
	}
}

func (i *Info) ReceiveGo(data string) {}

func (i *Info) ReceiveBack(source string, valSet []string) {}

/**** Send Handlers ****/

func (i *Info) SendPosition(dest string) error {
	connOut, err := net.DialTimeout("tcp", i.IP + ":" + dest, time.Duration(10) * time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send position to %s:%s \n", i.IP, dest)
			return err
		}
	}

	PositionMsg := Payload.Payload{ Source: i.Port, Dest: dest, Intent: constants.IntentSendPosition, Message: i.Neighbours}

	if err := json.NewEncoder(connOut).Encode(&PositionMsg); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", PositionMsg)
		return err
	}
	return nil
}

func (i *Info) SendGo(data, dest string) {}

func (i *Info) SendBack() {}

/**** Sanity Check Handler *****/
func (i *Info) SendPing(dest string) bool {
	connOut, err := net.DialTimeout("tcp", i.IP + ":" + dest, time.Duration(10) * time.Second)
	if err != nil {
		if _, ok := err.(net.Error); ok {
			fmt.Printf("Couldn't send ping to %s:%s \n", i.IP, dest)
			return false
		}
	}

	Ping := Payload.Payload{ Source: i.Port, Dest: dest, Intent: constants.IntentPing }
	if err := json.NewEncoder(connOut).Encode(&Ping); err != nil {
		fmt.Printf("Couldn't enncode message %v \n", Ping)
		return false
	}
	return true
}

/**** Node Communication Radar ****/
func (i *Info) ListenOnPort(){
	ln, err := net.Listen("tcp", fmt.Sprint(":" + i.Port))
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

		var reqMsg Payload.Payload
		if err := json.NewDecoder(connIn).Decode(&reqMsg); err != nil {
			fmt.Printf("Error decoding %v\n", err)
		}

		switch reqMsg.Intent {
		case constants.IntentSendPosition:
			i.ReceivePosition(reqMsg.Source, reqMsg.Message)
		case constants.IntentPing:
			fmt.Printf("Ping!\n")
		}
	}
}

