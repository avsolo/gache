package server

import (
	"fmt"
	"github.com/avsolo/gache/lib"
	s "github.com/avsolo/gache/storage"
)

var log *lib.Logger
var Store *s.Storage
var _s = fmt.Sprintf

func init() {
	log = lib.NewLogger("server")
	Store = s.NewStorage()
}
