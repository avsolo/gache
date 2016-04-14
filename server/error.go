package server

import "errors"


// ErrBadCommand returs when command (SET, GET...) not found in TCP request
var ErrBadCommand = errors.New("Bad command.")

// ErrBadArgs returns when in dict passed wrond key,value sequence
var ErrBadArgs = errors.New("Bad arguments.")

// ErrBadTTL returns when TTL not integer
var ErrBadTTL = errors.New("Bad TTL")

// ErrBadKey return if key is empty or has bad Go string key syntax
var ErrBadKey = errors.New("Bad key")

// ErrBadValue return when value is multiline
var ErrBadValue = errors.New("Bad value")

// ErrBadRequest return if request not recognized by patterns listed in routes
var ErrBadRequest = errors.New("Bad request.")
