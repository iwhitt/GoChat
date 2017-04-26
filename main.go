package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Println("pass c or s to run as server or client.")
		return
	}
	if os.Args[1] == "s" {
		RunServer()
	} else if os.Args[1] == "c" {
		RunClient()
	} else {
		fmt.Println("argument 1 was not \"s\" or \"c\"... exiting.")
	}

}
