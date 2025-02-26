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
		// fmt.Println("Error reading type. Error: ", err)
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

// get the Value properties and convert them to byte

func (v Value) marshall() []byte {
	switch v.typ {
	case "array":
		return v.marshallArray()
	case "bulk":
		return v.marshallBulk()
	case "string":
		return v.marshallString()
	case "null":
		return v.marshallNull()
	case "error":
		return v.marshallError()
	default:
		fmt.Println("Invalid type.")
		return []byte{}
	}
}

func (v Value) marshallArray() []byte {
	r := []byte{}

	length := len(v.arr)

	r = append(r, ARRAY)
	r = append(r, strconv.Itoa(len(v.arr))...)
	r = append(r, '\r')
	r = append(r, '\n')

	for i := 0; i < length; i++ {
		r = append(r, v.arr[i].marshall()...)
	}

	return r
}

func (v Value) marshallString() []byte {
	r := []byte{}

	r = append(r, STRING)
	r = append(r, v.str...)
	r = append(r, '\r')
	r = append(r, '\n')

	return r
}

func (v Value) marshallError() []byte {
	r := []byte{}

	r = append(r, ERROR)
	r = append(r, v.str...)
	r = append(r, '\r')
	r = append(r, '\n')

	return r
}

func (v Value) marshallInteger() []byte {
	r := []byte{}

	r = append(r, INTEGER)
	r = append(r, []byte(strconv.Itoa(v.num))...)
	r = append(r, '\r')
	r = append(r, '\n')

	return r
}

func (v Value) marshallBulk() []byte {
	r := []byte{}

	r = append(r, BULK)
	r = append(r, strconv.Itoa(len(v.bulk))...)
	r = append(r, '\r')
	r = append(r, '\n')
	r = append(r, []byte(v.bulk)...)
	r = append(r, '\r')
	r = append(r, '\n')

	return r
}

func (v Value) marshallNull() []byte {
	return []byte("$-1\r\n")
}

type Writer struct {
	writer io.Writer
}

func newWriter(writer io.Writer) *Writer {
	return &Writer{writer: writer}
}

func (w *Writer) write(v Value) error {
	bytes := v.marshall()

	_, err := w.writer.Write(bytes)
	if err != nil {
		fmt.Println("Error writing to client")
		return err
	}
	return nil
}
