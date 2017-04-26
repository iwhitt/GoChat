package main

type Client struct {
	co *net.TCPConn
	cc chan string
	cr *bufio.Reader
	cw *bufio.Writer
}

func main() {
	l, err := net.ListenTCP("TCP", ":8765")
	var cl []*Client
	for {
		conn, err = l.AcceptTCP()
		if err!=nil {
			fmt.Println("Error accepting connection.")
			continue
		}
		c := Client{conn,
			make(chan string),
			bufio.NewReader(conn),
			bufio.NewWriter(conn)}
		cl = append(cl, c)
		go clientListen(c)

	}
}

func clientListen(c Client) {
	for {
		str, err := c.cr.ReadString('\n')
		if err!= nil {
			fmt.Println("Error reading from client.")
			break
		}
		c.cc <- str
	}
}