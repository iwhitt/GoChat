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

	//r := bufio.NewReader(os.Stdin)
	// port := ":8765"
	// for {
	// 	fmt.Print("Input IP address: ")
	// 	str, _ := r.ReadString('\n')
	// 	str = strings.Trim(str, "\n")
	// 	str = strings.Trim(str, "\r")
	// 	addy := str + port
	// 	fmt.Println(addy)
	uname = "Default"
	addy := "127.0.0.1:8765"
	var err error
	addr, err = net.ResolveTCPAddr("tcp", addy)
	if err != nil {
		fmt.Print(err, "\t")
		fmt.Println("Error resolving address.")

	} // else {
	// 	fmt.Print("Input user name: ")
	// 	str, _ = r.ReadString('\n')
	// 	uname = strings.Trim(str, "\r")
	// 	break
	// }
	//}
	return addr, uname
}

func RunClient() {
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

	//go interruptListen(sc)
	go localListen(ic)
	go netListen(nc, r)

	w.WriteString("0" + username + "\n")
	w.Flush()

	for run {
		select {
		case _ = <-sc:
			w.WriteString("2")
			w.Flush()
			return
		case str := <-nc:
			str = strings.Trim(str, "\n")
			fmt.Println(str)
		case str := <-ic:
			if str == "quit" {
				fmt.Println("Quitting.")
				return
			}
			str = SanitizeString(str)
			w.WriteString("1" + str)
			w.Flush()
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
		str = SanitizeString(str)
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
