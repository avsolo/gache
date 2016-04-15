package server

import (
	"fmt"
    "strings"
	s "github.com/avsolo/gache/storage"
)

// Response is simple response object
type Response struct {
	Code int
	Error error
	Body string
}

// NewResponse create Response object frow s body and e error and returns
// pointer to it
func NewResponse(s string, e error) *Response {
	return &Response{Body:s, Error:e}
}

// Below list of routes

func routeSet(r *Request) *Response {
	err := Store.Set(r.Key, r.Value, r.TTL)
	if err != nil {
		return NewResponse("", err)
	}
	return NewResponse("[201]", nil)
}

func routeGet(r *Request) *Response {
	val, err := Store.Get(r.Key)
	if err != nil { return NewResponse("", s.ErrNotFound) }
	return NewResponse(fmt.Sprintf("%s", val), nil)
}

func routeUpdate(r *Request) *Response {
	err := Store.Update(r.Key, r.Value, r.TTL)
	if err != nil { return NewResponse("", err) }
	return NewResponse("[204]", nil)
}

func routeDelete(r *Request) *Response {
	Store.Delete(r.Key)
	return NewResponse("[204]", nil)
}

func routeService(r *Request) *Response {
    switch r.Key {
    case "flush":
        Store.Flush()
        return NewResponse("[204]", nil)
    default:
        return NewResponse("", ErrBadValue)
    }
}

// List routes
func routeLSet(r *Request) *Response {
    vals := strings.Split(r.Value, " ")
    rVal := []interface{}{}
    for _, v := range vals {
        rVal = append(rVal, v)
    }
    rVal = append(rVal, r.TTL)
	err := Store.LSet(r.Key, rVal...)
	if err != nil {
		return NewResponse("", err)
	}
	return NewResponse("[201]", nil)
}

func routeLPush(r *Request) *Response {
	err := Store.LPush(r.Key, r.Value)
	if err != nil { return NewResponse("", s.ErrNotFound) }
	return NewResponse("[204]", nil)
}

func routeLPop(r *Request) *Response {
	val, err := Store.LPop(r.Key)
	if err != nil { return NewResponse("", s.ErrNotFound) }
	return NewResponse(fmt.Sprintf("%s", val), nil)
}

// Dict routes
func routeDSet(r *Request) *Response {
    vals := strings.Split(r.Value, " ")

    rVal := []interface{}{}
    for _, v := range vals {
        rVal = append(rVal, v)
    }
    rVal = append(rVal, r.TTL)

	err := Store.DSet(r.Key, rVal...)
	if err != nil {
		return NewResponse("[400]", err)
	}
	return NewResponse("[201]", nil)
}

func routeDGet(r *Request) *Response {
	val, err := Store.DGet(r.Key, r.Value)
	if err != nil { return NewResponse("[404]", s.ErrNotFound) }
	return NewResponse(fmt.Sprintf("%s", val), nil)
}

func routeDAdd(r *Request) *Response {
    v := strings.Split(r.Value, " ")
	err := Store.DAdd(r.Key, strings.TrimSpace(v[0]), strings.TrimSpace(v[1]))
	if err != nil { return NewResponse("", s.ErrNotFound) }
	return NewResponse("[204]", nil)
}
