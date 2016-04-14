package server

import (
	"fmt"
	"avsolo/gache/lib"
	s "avsolo/gache/storage"
)

var log *lib.Logger
var Store *s.Storage
var _s = fmt.Sprintf

func init() {
	log = lib.NewLogger("server")
	Store = s.NewStorage()
}
