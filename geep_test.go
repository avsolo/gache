package main

import (
	"fmt"
	"github.com/avsolo/gache/lib"
	s "github.com/avsolo/gache/server"
)

var log *lib.Logger

func init() {
    log = lib.NewLogger("geep_test")
}

func main() {
	addr := lib.CliParams.ServerAddr

	s.ListenTCP(addr)

	g1 := s.NewGeepClient(addr)
	g2 := s.NewGeepClient(addr)

	go do(g1, 100, "SET k%d somevalue 100")
	go do(g2, 100, "GET k%d")

	select {}
	log.Debug("End select\n")
}

func do(g *s.GeepClient, count int, s string) {
	for i := 0; i < count; i++ {
		log.Debugf("Start %d\n", i)

		res, err := g.Send(fmt.Sprintf(s + "\r\n", i))
		if err != nil {
			log.Debugf("Send error: %s\n", err.Error())
			continue
		}
		log.Debugf("Send success: %s\n", res)
	}
	log.Debug("End cycle\n")
}

