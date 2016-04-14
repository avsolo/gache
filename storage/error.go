package storage

import (
	"errors"
	"avsolo/gache/lib"
)

var log *lib.Logger

func init() {
	log = lib.NewLogger("storage")
}

// ErrAlreadyExists uses when client trying to set existing key
// By default client allowed to update same key, not to set it
var ErrAlreadyExists = errors.New("Key already exists")

// ErrNotFound uses when client trying to update non existing key
var ErrNotFound = errors.New("Key not found")

var ErrBadTTL = errors.New("Bad TTL")

// ErrNotList represent error when user trying pop or push not list item
var ErrNotList = errors.New("Key not list")

// ErrEmpty represent error when user trying get sub-element of any container
var ErrEmpty = errors.New("Contaiter empty")

var ErrNoExpire = errors.New("Key has no expire")

var ErrBadMap = errors.New("Bad argument(s) for hash")

var ErrNotDict = errors.New("Key not dict")
