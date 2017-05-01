package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type Client struct {
	co *net.TCPConn
	//cc chan string
	cr   *bufio.Reader
	cw   *bufio.Writer
	name string
	room int
}

type message struct {
	roomID  int
	message string
}

var rooms []int

func RunServer() {
	run := true

	a, err := net.ResolveTCPAddr("tcp", ":8765")
	if err != nil {
		fmt.Println("Error resolving address.")
		run = false
	}

	l, _ := net.ListenTCP("tcp", a)
	if err != nil {
		fmt.Println("Error listening.")
		run = false
	}

	var cl []*Client

	ic := make(chan string)       // Input channel
	cc := make(chan message)      // Client channel
	jc := make(chan *net.TCPConn) // Join channel
	lc := make(chan *Client)      // Leave channel

	go inputListen(ic)
	go joinListen(l, jc)

	for run {
		select {
		case cmd := <-ic:
			if cmd == "quit" {
				fmt.Println("Server quitting.")
				return
			}
		case conn := <-jc:
			cliRoomID := -1
			for i := 0; i < len(rooms); i++ {
				if rooms[i] < 2 {
					cliRoomID = i
					rooms[i]++
					break
				}
			}
			if cliRoomID == -1 {
				rooms = append(rooms, 1)
				cliRoomID = len(rooms) - 1
			}
			fmt.Println("new client going to room", cliRoomID)
			c := Client{conn,
				bufio.NewReader(conn),
				bufio.NewWriter(conn),
				"NONAME", cliRoomID}
			cl = append(cl, &c)
			go clientListen(c, cc, lc)
		case msg := <-cc:
			for _, c := range cl {
				if c.room == msg.roomID {
					c.cw.WriteString(msg.message)
					c.cw.Flush()
				}
			}
			str := strings.Trim(msg.message, "\n")
			fmt.Println(str)
		case c := <-lc:
			for i, tc := range cl {
				if c == tc {
					fmt.Println(rooms[tc.room])
					rooms[tc.room]--
					fmt.Println(rooms[tc.room])
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
		if err != nil {
			fmt.Println("Error reading from console.")
			continue
		}
		str = strings.Trim(str, "\r")
		fmt.Println(str)
		if str == "quit" {
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
		if err != nil {
			fmt.Print("Error accepting connection.\r")
			continue
		}
		jc <- conn
		fmt.Println("New client joined.")
	}
}

func clientListen(c Client, cc chan message, lc chan *Client) {
	for {
		str, err := c.cr.ReadString('\n')
		if err != nil {
			rooms[c.room]--
			fmt.Println("Error reading from client. Dropping client", c.name)
			lc <- &c
			break
		}

		tag := string(str[0])
		mes := string(str[1:])
		msg := message{c.room, mes}
		if tag == "0" {
			mes = strings.Trim(mes, "\n")
			fmt.Println(c.name, "is now", mes)
			c.name = mes
		} else if tag == "1" {
			msg.message = c.name + ": " + msg.message
			cc <- msg
		} else if tag == "2" {
			fmt.Println("Client", c.name, "quit.")
			lc <- &c
			break
		}
	}
}
