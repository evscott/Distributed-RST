package main

/* Al useful imports */
import (
	"bufio"
	"flag"
	"fmt"
	"github.com/evscott/DistroA1/constants"
	"github.com/evscott/DistroA1/models/Node"
	"net"
	"os"
	"strings"
	"time"
)

func main(){
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
	time.Sleep(time.Second/2)
	exposeCLI(node)
}

func exposeCLI(node *Node.Info) {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf(constants.StdIn, "enter cmd")
		if rawUserInput, err := reader.ReadString('\n'); err == nil {
			userInput := strings.Split(rawUserInput, "\n")[0]

			switch userInput {
			case "start":
				node.Start()
			case "ping":
				fmt.Printf(constants.StdIn, "enter port")
				if rawUserInput, err := reader.ReadString('\n'); err == nil {
					userInput = strings.Split(rawUserInput, "\n")[0]
					node.SendPing(userInput)
				}
			case "show neighbors":
				fmt.Printf(constants.StdOut, node.Neighbours, "")
			case "show proc_known":
				fmt.Printf(constants.StdOut, node.ProcKnown, "")
			case "show channels_known":
				fmt.Printf(">>")
				for _, p := range node.ChannelsKnown {
					fmt.Printf(" (%v, %v)", p.I, p.J)
				}
				fmt.Printf("\n")
			case "exit":
				fmt.Printf(constants.StdOut, "shutting down node", node.Port)
				os.Exit(0)
			default:
				fmt.Printf(constants.StdOut, "no command found for", userInput)
			}
		}
	}
}