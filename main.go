package main

import (
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/evscott/DistroA1/Node"
)

// main is the entry point for this distributed system
//
func main() {
	ipAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Print(err)
		return
	}
	ip := strings.Split(ipAddr[0].String(), "/")[0]

	runExample(ip)
}

// runExample creates an exemplary distributed rooted spanning tree from a generic communication graph
func runExample(ip string) {
	n1 := Node.Create(ip, "8000", []string{"8001", "8006", "8007"})
	go n1.ListenOnPort()

	n2 := Node.Create(ip, "8001", []string{"8000", "8002", "8007"})
	go n2.ListenOnPort()

	n3 := Node.Create(ip, "8002", []string{"8001", "8003"})
	go n3.ListenOnPort()

	n4 := Node.Create(ip, "8003", []string{"8002", "8007", "8004", "8005"})
	go n4.ListenOnPort()

	n5 := Node.Create(ip, "8004", []string{"8005", "8003"})
	go n5.ListenOnPort()

	n6 := Node.Create(ip, "8005", []string{"8003", "8004", "8006"})
	go n6.ListenOnPort()

	n7 := Node.Create(ip, "8006", []string{"8000", "8007", "8003", "8005"})
	go n7.ListenOnPort()

	n8 := Node.Create(ip, "8007", []string{"8007", "8000", "8001", "8003", "8006"})
	go n8.ListenOnPort()

	time.Sleep(time.Second / 10)

	n1.StartRST()

	time.Sleep(time.Second)
}