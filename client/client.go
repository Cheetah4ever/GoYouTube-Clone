package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	fmt.Println("Client running on Local Host: 4200")
	//conn, err := net.Dial("udp4", "localhost:1053")
	conn, err := net.Dial("tcp", "localhost:4200")
	check(err)
	defer conn.Close()

	for {
		send := bufio.NewScanner(os.Stdin)
		if send.Scan() {
			conn.Write(send.Bytes())
		}
		recv := make([]byte, 4048)
		n, error := conn.Read(recv)
		if error != nil {
			fmt.Println(error)
		}
		fmt.Fprintf(os.Stdout, "%d\n", recv[:n])
	}
}
