package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Value struct {
	typ  string
	str  string
	num  int
	bulk string
	arr  []Value
}

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

type Resp struct {
	reader *bufio.Reader
}

func newResp(reader io.Reader) *Resp {
	return &Resp{reader: bufio.NewReader(reader)}
}

func (rd *Resp) read() (Value, error) {
	typ, err := rd.reader.ReadByte()

	if err != nil {
		fmt.Println("Error reading type.")
		return Value{}, err
	}

	switch typ {
	case ARRAY:
		return rd.readArray()
	case BULK:
		return rd.readBulk()
	default:
		return Value{}, nil
	}
}

func (rd *Resp) readArray() (Value, error) {
	val := Value{}
	val.typ = "array"

	length, _, err := rd.readInteger()
	if err != nil {
		fmt.Println("Error reading integer.")
		return val, err
	}

	val.arr = make([]Value, length)
	for i := 0; i < length; i++ {
		val1, err := rd.read()
		if err != nil {
			fmt.Println("Error reading value.")
			return val, err
		}

		val.arr[i] = val1
	}

	return val, err

}

func (rd *Resp) readBulk() (Value, error) {
	val := Value{}
	val.typ = "bulk"

	length, _, err := rd.readInteger()
	if err != nil {
		fmt.Println("Error reading integer.")
		return val, err
	}

	bulk := make([]byte, length)
	rd.reader.Read(bulk)

	rd.reader.ReadLine()

	val.bulk = string(bulk)

	return val, nil
}

func (rd *Resp) readLine() (line []byte, n int, err error) {

	for {
		value, err := rd.reader.ReadByte()

		if err != nil {
			fmt.Println("Error reading line.")
			return nil, 0, err
		}

		n += 1
		line = append(line, value)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}

	}

	return line[:len(line)-2], n, nil
}

func (rd *Resp) readInteger() (int, int, error) {
	line, n, err := rd.readLine()

	if err != nil {
		fmt.Println("Error reading line.")
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		fmt.Println("Error parsing integer.")
		return 0, n, err
	}

	return int(i64), n, nil
}
