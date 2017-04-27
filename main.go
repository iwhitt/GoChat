package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
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

// SanitizeString cleans input so it can be safely written to connection.
func SanitizeString(input string) string {
	result := input
	if runtime.GOOS == "windows" {
		result = strings.TrimRight(input, "\r\n")
		result = result + "\n"
	}
	return result
}
