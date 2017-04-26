package main

import (
	"net"
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Client struct {
	co *net.TCPConn
	//cc chan string
	cr *bufio.Reader
	cw *bufio.Writer
	name string
}

func main() {
	run := true

	a, err := net.ResolveTCPAddr("tcp", ":8765")
	if err!=nil {
		fmt.Println("Error resolving address.")
		run = false
	}

	l, _ := net.ListenTCP("tcp",a)
	if err!=nil {
		fmt.Println("Error listening.")
		run = false
	}

	var cl []*Client

	ic := make(chan string)		// Input channel
	cc := make(chan string)		// Client channel
	jc := make(chan *net.TCPConn)	// Join channel
	lc := make(chan *Client)	// Leave channel

	go inputListen(ic)
	go joinListen(l, jc)

	for run{
		select {
		case cmd := <-ic:
			if cmd=="quit" {
				return
			}
		case conn := <-jc:
			c := Client{conn,
				    bufio.NewReader(conn),
				    bufio.NewWriter(conn),
				    "NONAME"}
			cl = append(cl, &c)
			go clientListen(c, cc, lc)
		case str := <-cc:
			for _, c := range cl {
				c.cw.WriteString(str)
				c.cw.Flush()
			}
			fmt.Println(str)
		case c := <- lc:
			for i, tc := range cl {
				if c==tc {
					cl = append(cl[:i], cl[i+1:]...)
					fmt.Println("Client removed.")
				}
			}
		}
	}
}

func inputListen(ic chan string) {
	fmt.Println("Reading console input.")
	r := bufio.NewReader(os.Stdin)
	for {
		str, err := r.ReadString('\n')
		if err!= nil {
			fmt.Println("Error reading from console.")
			continue
		}
		str = strings.TrimRight(str,"\n")
		fmt.Println(str)
		if str=="quit" {
			fmt.Println("QUIT")
			ic <- str
			break
		}
		ic <- str
	}
}

func joinListen(l *net.TCPListener, jc chan *net.TCPConn) {
	fmt.Println("Listening...")
	for {
		conn, err := l.AcceptTCP()
		if err!=nil {
			fmt.Print("Error accepting connection.\r")
			continue
		}
		jc <- conn
		fmt.Println("New client joined.")
	}
}

func clientListen(c Client, cc chan string, lc chan *Client) {
	for {
		str, err := c.cr.ReadString('\n')
		if err!= nil {
			fmt.Println("Error reading from client. Dropping client.")
			lc <- &c
			break
		}
		tag := string(str[0])
		mes := string(str[1:])
		if tag=="0" {
			c.name = mes
		} else if tag=="1" {
			mes = c.name + ": " + mes
			cc <- mes
		} else if tag=="2" {
			fmt.Println("Client",c.name,"quit.")
			lc <- &c
			break
		}
	}
}