package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"strings"
	"time"

	"github.com/evscott/DistroA1/Node"
)

/* TODO */
func main() {
	port := flag.String("port", "", "The port to run this node on")
	flag.Parse()
	neighbors := flag.Args()

	ipAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Print(err)
		return
	}
	ip := strings.Split(ipAddr[0].String(), "/")[0]

	node := Node.Create(ip, *port, neighbors)

	go node.ListenOnPort()
	time.Sleep(time.Second / 10)
	runCLI(node)
}

/* TODO */
func runCLI(node *Node.Info) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(">> enter cmd: ")
		if rawUserInput, err := reader.ReadString('\n'); err == nil {
			userInput := strings.Split(rawUserInput, "\n")[0]

			switch userInput {
			case "startCG":
				node.StartCG()
			case "startRST":
				node.StartRST()
			case "ping":
				fmt.Printf(">>>> enter port: ")
				if rawUserInput, err := reader.ReadString('\n'); err == nil {
					userInput = strings.Split(rawUserInput, "\n")[0]
					node.SendPing(userInput)
				}
			case "show neighbors":
				fmt.Printf(">> [%v]", node.Neighbours)
			case "show proc_known":
				for key, _ := range node.ProcKnown {
					fmt.Printf(">>{%s}\n", key)
				}
			case "show channels_known":
				for key, _ := range node.ChannelsKnown {
					fmt.Printf(">>(%s, %s)\n", key.I, key.J)
				}
			case "show inbox":
				for key, _ := range node.MessageInbox {
					fmt.Printf(">>[%v]\n", key)
				}
			case "show children":
				for key, _ := range node.Children {
					fmt.Printf(">>[%v]\n", key)
				}
			case "show parent":
				fmt.Printf(">> %s", node.Parent)
			case "exit":
				fmt.Printf(">> shutting down node %v\n", node.Port)
				os.Exit(0)
			default:
				fmt.Printf(">> no command found for %v\n", userInput)
			}
		}
	}
}
