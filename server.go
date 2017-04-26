package main

import (
	"net"
	"bufio"
	"fmt"
)

type Client struct {
	co *net.TCPConn
	//cc chan string
	cr *bufio.Reader
	cw *bufio.Writer
}

func main() {
	a, _ := net.ResolveTCPAddr("tcp", ":8765")

	l, _ := net.ListenTCP("tcp",a)

	var cl []*Client

	cc := make(chan string)		// Client channel
	jc := make(chan *net.TCPConn)	// Join channel

	go joinListen(l, jc)

	for {
		select {
		case conn := <- jc:
			c := Client{conn,
				    bufio.NewReader(conn),
				    bufio.NewWriter(conn)}
			cl = append(cl, &c)
			go clientListen(c, cc)
		case str := <- cc:
			for _, c := range cl {
				c.cw.WriteString(str)
				c.cw.Flush()
			}
			fmt.Println(str)
		}
	}
}

func joinListen(l *net.TCPListener, jc chan *net.TCPConn) {
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

func clientListen(c Client, cc chan string) {
	for {
		str, err := c.cr.ReadString('\n')
		if err!= nil {
			fmt.Println("Error reading from client. Dropping client.")
			break
		}
		cc <- str
	}
}