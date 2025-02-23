package main

import (
	"fmt"
	"sync"
)

const (
	PING    = "PING"
	SET     = "SET"
	GET     = "GET"
	HSET    = "HSET"
	HGET    = "HGET"
	HGETALL = "HGETALL"
)

var STORAGE = map[string]string{}

// ensure no concurency issues
var STORAGE_LOCK = sync.RWMutex{}

var handlers = map[string]func([]Value) Value{
	PING: ping,
	SET:  set,
	GET:  get,
}

func ping(args []Value) Value {
	if len(args) != 0 {
		return Value{typ: "string", str: args[0].bulk}
	}
	return Value{typ: "string", str: "PONG"}
}

func set(args []Value) Value {
	key := args[0].bulk

	// check if key already exists
	_, ok := STORAGE[key]
	if ok {
		return Value{typ: "string", str: "ERROR key already exists."}
	}

	STORAGE_LOCK.Lock()
	STORAGE[key] = args[1].bulk
	STORAGE_LOCK.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) < 1 {
		fmt.Println("Invalid number of arguments.")
		return Value{typ: "error", str: "Invalid number of arguments."}
	}

	key := args[0].bulk

	STORAGE_LOCK.RLock()
	value, ok := STORAGE[key]
	STORAGE_LOCK.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}
