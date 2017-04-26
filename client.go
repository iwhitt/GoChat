package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
)

func setup() (*net.TCPAddr, string) {
	var addr *net.TCPAddr
	var uname string
	port := ":8765"
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Input IP address: ")
		str, _ := r.ReadString('\n')
		fmt.Println(str)
		str = strings.Trim(str, "\n")
		str = strings.Trim(str, "\r")
		fmt.Println(str)
		addy := str + port
		fmt.Println(addy)
		var err error
		addr, err = net.ResolveTCPAddr("tcp", addy)
		if err != nil {
			fmt.Print(err, "\t")
			fmt.Println("Error resolving address.")
		} else {
			fmt.Print("Input user name: ")
			str, _ = r.ReadString('\n')
			uname = strings.TrimRight(str, "\n")
			break
		}
	}
	return addr, uname
}

func main() {
	a, username := setup()

	run := true

	c, err := net.DialTCP("tcp", nil, a)
	if err != nil {
		fmt.Println("Error connecting.")
		run = false
	}

	ic := make(chan string)    // Input channel
	nc := make(chan string)    // Net channel
	sc := make(chan os.Signal) // Signal channel

	r := bufio.NewReader(c)
	w := bufio.NewWriter(c)
	w.WriteString("0" + username)
	w.Flush()

	go interruptListen(sc)
	go localListen(ic)
	go netListen(nc, r)

	for run {
		select {
		case _ = <-sc:
			w.WriteString("2")
			w.Flush()
			return
		case str := <-nc:
			fmt.Println(str)
		case str := <-ic:
			w.WriteString("1" + str)
			w.Flush()
			fmt.Println("?")
		}
	}
}

func interruptListen(ic chan os.Signal) {
	for {
		signal.Notify(ic, os.Interrupt)
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
		str = strings.TrimRight(str, "\n")
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
