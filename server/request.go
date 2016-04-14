// Package server provides ...
package server

import (
	"errors"
    "regexp"
	"strings"
	"strconv"
)

// All our commands available by TCP must be listed in this list
// and in var pathes.
const (
	CMD_SET	  = "SET"
	CMD_GET	  = "GET"
	CMD_UPD	  = "UPD"
	CMD_DEL	  = "DEL"

	CMD_LSET  = "LSET"
	CMD_LPUSH = "LPUSH"
	CMD_LPOP  = "LPOP"

	CMD_DSET  = "DSET"
	CMD_DGET  = "DGET"
	CMD_DADD  = "DADD"
	CMD_DDEL  = "DDEL"

	CMD_OPT  = "OPT"
)

// path keep compiled regexp and callable func for routing
type path struct {
    Re *regexp.Regexp
    Method func(r *Request) *Response
}

// Set and Get patterns
var rmc = regexp.MustCompile // Just shorcut
var setPtn = rmc(`^(?P<key>\w+)\s+(?P<value>.*)\s+(?P<ttl>\d+)$`)
var dAddPtn = rmc(`^(?P<key>\w+)\s+(?P<value>.*)$`)
var getPtn = rmc(`^(?P<key>\w+)$`)

// List of routes
var pathes = map[string]*path{
    CMD_SET: &path{setPtn, routeSet},
    CMD_GET: &path{getPtn, routeGet},
    CMD_UPD: &path{setPtn, routeUpdate},
    CMD_DEL: &path{getPtn, routeDelete},

    CMD_LSET: &path{setPtn, routeLSet},
    CMD_LPUSH: &path{dAddPtn, routeLPush},
    CMD_LPOP: &path{getPtn, routeLPop},

    CMD_DSET: &path{setPtn, routeDSet},
    CMD_DGET: &path{dAddPtn, routeDGet},
    CMD_DADD: &path{dAddPtn, routeDAdd},

    CMD_OPT: &path{getPtn, routeService},
}

// Request provide general request object and keep all required data
type Request struct {
	Cmd string
	Key string
	Value string
	TTL int
	Method func(r *Request) *Response
    Raw map[string]string
}

// NewRequest get sting, split and do base validation (number of params,
// type of value, etc.)
func NewRequest(in string) (*Request, error) {
	in = strings.TrimSpace(in)
	if in == "" {
		log.Warn("Empty request")
		return nil, ErrBadRequest
	}
	fp := strings.SplitN(in, " ", 2) // First, get CMD name
	if len(fp) != 2 {
		log.Warnf("Request creating error: Request: %s", in)
		return nil, ErrBadRequest
	}

	// Create init params
	r := &Request{Cmd: strings.TrimSpace(fp[0]), Raw: map[string]string{}}
	if _, found := pathes[r.Cmd]; !found {
		log.Warnf("Cmd '%s' unknown", r.Cmd)
		return nil, ErrBadCommand
	}

	// Match reqeust string
	n1 := pathes[r.Cmd].Re.SubexpNames()
	args := pathes[r.Cmd].Re.FindStringSubmatch(fp[1])
	if args == nil {
        emsg := _s("Error matching args from string: %#v", fp[1])
        log.Warnf(emsg)
		return nil, errors.New(emsg)
	}

	// Clean params
	var err error
	for i, n := range args {
		switch n1[i] {
		case "key":
			r.Key = n
		case "value":
			r.Value = n
		case "ttl":
			if r.TTL, err = strconv.Atoi(n); err != nil {
				return nil, ErrBadTTL
			}
		case "":
			continue
		default:
			return nil, errors.New(_s("Unknown key '%s'", n1[i]))
		}
	}
	r.Method = pathes[r.Cmd].Method
	return r, nil
}

func (r *Request) Route() *Response {
	return r.Method(r)
}
