package main

import (
	"fmt"
	"log"
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

var HSTORAGE = map[string]map[string]string{}
var HSTORAGE_LOCK = sync.RWMutex{}

var handlers = map[string]func([]Value) Value{
	PING:    ping,
	SET:     set,
	GET:     get,
	HSET:    hSet,
	HGET:    hGet,
	HGETALL: hGetAll,
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
	STORAGE_LOCK.Lock()
	_, ok := STORAGE[key]

	if ok {
		return Value{typ: "string", str: "ERROR key already exists."}
	}

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

func hGetAll(args []Value) Value {
	if len(args) < 1 {
		return Value{typ: "error", str: "Invalid number of arguments for an hgetall command."}
	}

	key := args[0].bulk

	HSTORAGE_LOCK.RLock()
	value, ok := HSTORAGE[key]
	HSTORAGE_LOCK.RUnlock()

	if !ok {
		fmt.Println("Requested key does not exist.")
		return Value{typ: "error", str: "Requested key does not exist."}
	}

	values := []Value{}
	for k, v := range value {
		values = append(values, Value{typ: "bulk", bulk: k})
		values = append(values, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", arr: values}
}

func hSet(args []Value) Value {
	if len(args) < 3 {
		return Value{typ: "error", str: "Invalid number of arguments for an hset command."}
	}

	key := args[0].bulk
	// field := args[1].bulk
	// value := args[2].bulk

	if _, ok := HSTORAGE[key]; !ok {
		HSTORAGE[key] = map[string]string{}
	}

	for i := 1; i < len(args)-1; i++ {
		k := args[i].bulk
		v := args[i+1].bulk

		HSTORAGE_LOCK.Lock()
		HSTORAGE[key][k] = v
		HSTORAGE_LOCK.Unlock()
	}

	return Value{typ: "string", str: "OK"}
}

func hGet(args []Value) Value {
	if len(args) < 2 {
		return Value{typ: "error", str: "Invalid number of arguments for an hget command."}
	}

	key := args[0].bulk
	field := args[1].bulk

	HSTORAGE_LOCK.RLock()
	log.Printf("Check for %s key and %s field.", key, field)
	value, ok := HSTORAGE[key][field]
	if !ok {
		return Value{typ: "error", str: "Requested field does not exist."}
	}

	HSTORAGE_LOCK.RUnlock()
	return Value{typ: "bulk", bulk: value}
}
