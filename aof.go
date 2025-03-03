package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sync"
	"time"
)

type AOF struct {
	file *os.File
	rd   *bufio.Reader
	mut  sync.Mutex
}

func (aof *AOF) Write(value Value) error {
	aof.mut.Lock()
	defer aof.mut.Unlock()

	_, err := aof.file.Write(value.marshall())
	if err != nil {
		fmt.Printf("Error writing to file. %v", err)
		return err
	}

	return nil
}

func newAof(path string) (*AOF, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error opening file.")
		return nil, err
	}

	aof := &AOF{
		file: file,
		rd:   bufio.NewReader(file),
	}

	go func() {
		for {
			aof.mut.Lock()

			aof.file.Sync()

			aof.mut.Unlock()

			time.Sleep(1 * time.Second)
		}
	}()

	return aof, nil
}

func (aof *AOF) close() error {
	aof.mut.Lock()
	defer aof.mut.Unlock()

	err := aof.file.Close()
	if err != nil {
		fmt.Println("Error closing file.")
		return err
	}

	return nil
}

func (aof *AOF) read(callback func(value Value)) error {
	aof.mut.Lock()
	defer aof.mut.Unlock()

	resp := newResp(aof.file)

	for {
		val, err := resp.read()

		if err == nil {
			callback(val)
		}

		if err == io.EOF {
			break
		}

	}

	return nil
}
