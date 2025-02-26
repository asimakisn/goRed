package main

import (
	"fmt"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	l, err := net.Listen("tcp", ":6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	aof, err := newAof("testAof.aof")
	if err != nil {
		fmt.Println(err)
	}

	defer aof.file.Close()

	aof.read(func(val Value) {
		command := strings.ToUpper(val.arr[0].bulk)
		args := val.arr[1:]

		handler, ok := handlers[command]
		if !ok {
			fmt.Println("Invalid command:", command)
			return
		}

		handler(args)
	})

	conn, err := l.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer conn.Close()

	for {

		rsp := newResp(conn)

		val, err := rsp.read()
		if err != nil {
			fmt.Println(err)
			return
		}

		if val.typ != "array" || len(val.arr) == 0 {
			fmt.Println("Invalid request, expected array.")
			continue
		}

		writer := newWriter(conn)

		command := strings.ToUpper(val.arr[0].bulk)

		if command == "COMMAND" {
			writer.write(Value{typ: "string", str: "OK"})
			continue
		}

		args := val.arr[1:]

		handler, ok := handlers[command]

		if command == "SET" || command == "HSET" {
			err := aof.Write(val)
			if err != nil {
				fmt.Println("Error writing to aof.")
			}
		}

		if !ok {
			fmt.Println("Invalid command:", command)
			errorMessage := "Invalid command: " + command
			writer.write(Value{typ: "error", str: errorMessage})
			continue
		}

		result := handler(args)

		writer.write(result)
	}
}
