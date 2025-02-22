package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {

		fmt.Println("Reading from client...")
		rsp := newResp(conn)

		val, err := rsp.read()
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Successfully read from client:", val)
		fmt.Println(val)

		conn.Write([]byte("+OK\r\n"))
	}
}
