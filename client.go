package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	run := true

	a, err := net.ResolveTCPAddr("tcp", "206.21.94.53:8765")
	if err != nil {
		fmt.Print(err, "\t")
		fmt.Println("Error resolving address.")
		run = false
	}

	c, err := net.DialTCP("tcp", nil, a)
	if err != nil {
		fmt.Println("Error connecting.")
		run = false
	}

	ic := make(chan string) // Input channel
	nc := make(chan string) // Net channel

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)

	go localListen(ic)
	go netListen(nc, r)

	for run {
		select {
		case str := <-nc:
			fmt.Println(str)
		case str := <-ic:
			if str == "quit" {
				run = false
				break
			}
			w.WriteString(str)
			w.Flush()
		}
	}
}

func localListen(ic chan string) {
	r := bufio.NewReader(os.Stdin)
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from console.")
			continue
		}
		if str == "quit" {
			ic <- str
			break
		}
		ic <- str
	}
}

func netListen(nc chan string, r *bufio.Reader) {
	for {
		str, err := r.ReadString('\n')
		if err != nil {
			fmt.Println("Error reading from connection.")
			break
		}
		nc <- str
	}
}
