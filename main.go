package main

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/evscott/Distributed-RST/Node"
)

// main is the entry point for this distributed system.
func main() {
	ipAddr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Print(err)
		return
	}
	ip := strings.Split(ipAddr[0].String(), "/")[0]

	runExample(ip)
}

// displayNeighbors prints to standard output the neighbors belonging to a set of nodes.
func displayNeighbors(nodeArr ...*Node.Info) {
	for _, node := range nodeArr {
		if port, err := strconv.Atoi(node.Port); err == nil {
			p := port % 1000
			fmt.Printf("Node %d: neighbors: {", p)
			for _, neighbour := range node.Neighbours {
				if c, err := strconv.Atoi(neighbour); err == nil {
					c := c % 1000
					fmt.Printf("%d,", c)
				}
			}
			fmt.Printf("}\n")
		}
	}
}

// displayNeighbors prints to standard output the children belonging to a set of nodes.
func displayChildren(nodeArr ...*Node.Info) {
	for _, node := range nodeArr {
		if port, err := strconv.Atoi(node.Port); err == nil {
			p := port % 1000
			fmt.Printf("Node %d: children: {", p)
			for child, _ := range node.Children {
				if c, err := strconv.Atoi(child); err == nil {
					c := c % 1000
					fmt.Printf("%d,", c)
				}
			}
			fmt.Printf("}\n")
		}
	}
}

// displayRSTAdjMatrix prints to standard output the adjacency matrix
// representation of a rooted spanning tree given a set of nodes.
func displayRSTAdjMatrix(nodeArr ...*Node.Info) {
	size := len(nodeArr)

	adjMatrix := make([][]uint8, size)
	for i := range adjMatrix {
		adjMatrix[i] = make([]uint8, size)
	}

	for _, node := range nodeArr {
		if port, err := strconv.Atoi(node.Port); err == nil {
			p := port % 1000

			for child, _ := range node.Children {
				if c, err := strconv.Atoi(child); err == nil {
					c := c % 1000

					adjMatrix[c][p] = 1
				}
			}
		}
	}

	printMatrix(size, adjMatrix)
}

// displayCGAdjMatrix prints to standard output the adjacency matrix
// representation of a communication graph given a set of nodes.
func displayCGAdjMatrix(nodeArr ...*Node.Info) {
	size := len(nodeArr)

	adjMatrix := make([][]uint8, size)
	for i := range adjMatrix {
		adjMatrix[i] = make([]uint8, size)
	}

	for _, node := range nodeArr {
		if port, err := strconv.Atoi(node.Port); err == nil {
			p := port % 1000

			for _, neighbour := range node.Neighbours {
				if c, err := strconv.Atoi(neighbour); err == nil {
					c := c % 1000

					adjMatrix[c][p] = 1
					adjMatrix[p][c] = 1
				}
			}
		}
	}

	printMatrix(size, adjMatrix)
}

func printMatrix(size int, adjMatrix [][]uint8) {
	for i := 0;  i < size; i++ {
		fmt.Printf("[")
		for j := 0; j < size; j++ {
			if j == size-1 {
				fmt.Printf("%d", adjMatrix[i][j])
			} else {
				fmt.Printf("%d,", adjMatrix[i][j])
			}
		}
		fmt.Printf("]\n")
	}
}

// runExample creates a rooted spanning tree from an arbitrary communication graph and display it.
//
// This example borrows the example put forward by Raynal in `Distributed Algorithm's for Message Passing Systems, but
// substitutes letters for numbers in the identification of nodes.
// `a` = `8000`
// `b` = `8001`
// `c` = `8002`
// ... etc
func runExample(ip string) {
	n1 := Node.Create(ip, "8000", []string{"8001", "8006", "8007"})
	n2 := Node.Create(ip, "8001", []string{"8000", "8002", "8007"})
	n3 := Node.Create(ip, "8002", []string{"8001", "8003"})
	n4 := Node.Create(ip, "8003", []string{"8002", "8007", "8004", "8005"})
	n5 := Node.Create(ip, "8004", []string{"8005", "8003"})
	n6 := Node.Create(ip, "8005", []string{"8003", "8004", "8006"})
	n7 := Node.Create(ip, "8006", []string{"8000", "8007", "8003", "8005"})
	n8 := Node.Create(ip, "8007", []string{"8007", "8000", "8001", "8003", "8006"})

	go n1.ListenOnPort()
	go n2.ListenOnPort()
	go n3.ListenOnPort()
	go n4.ListenOnPort()
	go n5.ListenOnPort()
	go n6.ListenOnPort()
	go n7.ListenOnPort()
	go n8.ListenOnPort()

	time.Sleep(time.Second / 10)

	fmt.Println()

	n1.Start()

	time.Sleep(time.Second)

	fmt.Printf("\nCommunication graph: \n\n")

	displayNeighbors(n1, n2, n3, n4, n5, n6, n7, n8)

	fmt.Println()

	displayCGAdjMatrix(n1, n2, n3, n4, n5, n6, n7, n7)

	fmt.Printf("\nRooted spanning tree: \n\n")

	displayChildren(n1, n2, n3, n4, n5, n6, n7, n8)

	fmt.Println()

	displayRSTAdjMatrix(n1, n2, n3, n4, n5, n6, n7, n8)
}