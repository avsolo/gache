package main

import (
	"fmt"
	"time"
	"testing"
	"github.com/avsolo/gache/lib"
	"github.com/avsolo/gache/server"
)

var _s = fmt.Sprintf
var log *lib.Logger
var addr = lib.CliParams.ServerAddr
var srv = server.NewServer(addr)
var cln *server.Client

func init() {
	log = lib.NewLogger("benchmark")
    go func() {
        srv.ListenTCP()
        log.Info("Server stopped")
    }()
    // Sometimest client trying connect before server start.
    // So, we just wait a bit
    time.Sleep(time.Duration(100 * time.Microsecond))
    cln = server.NewClient(addr)
}

func BenchmarkSetGet(b *testing.B) {
    cln.Send("OPT flush")
	for i := 0; i < b.N; i++ {
		_, err := cln.Send(fmt.Sprintf("SET k%d somevalue 100\r\n", i))
		if err != nil {
			panic("Unable make TCP reqeust. Error: " + err.Error())
		}
		_, err = cln.Send(fmt.Sprintf("GET k%d\r\n", i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkSetUpdate(b *testing.B) {
    cln.Send("OPT flush")
	for i := 0; i < b.N; i++ {
		_, err := cln.Send(fmt.Sprintf("SET k%d somevalue 100\r\n", i))
		if err != nil {
			panic("Unable make TCP reqeust. Error: " + err.Error())
		}
		_, err = cln.Send(fmt.Sprintf("UPD k%d val2 200\r\n", i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkLSetLPush(b *testing.B) {
    cln.Send("OPT flush")
	for i := 0; i < b.N; i++ {
		_, err := cln.Send(fmt.Sprintf("LSET k%d v1 v2 v3 100\r\n", i))
		if err != nil {
			panic("Unable make TCP reqeust. Error: " + err.Error())
		}
		_, err = cln.Send(fmt.Sprintf("LPUSH k%d v4\r\n", i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkLSetLPop(b *testing.B) {
    cln.Send("OPT flush")
	for i := 0; i < b.N; i++ {
		_, err := cln.Send(fmt.Sprintf("LSET k%d v1 v2 v3 100\r\n", i))
		if err != nil {
			panic("Unable make TCP reqeust. Error: " + err.Error())
		}
		_, err = cln.Send(fmt.Sprintf("LPOP k%d\r\n", i))
		if err != nil {
			panic(err.Error())
		}
	}
}

func BenchmarkDSetDAdd(b *testing.B) {
    cln.Send("OPT flush")
	for i := 0; i < b.N; i++ {
		_, err := cln.Send(fmt.Sprintf("DSET k%d k1,v1, k2,v2, k3,v3 100\r\n", i))
		if err != nil {
			panic("Unable make TCP reqeust. Error: " + err.Error())
		}
		_, err = cln.Send(fmt.Sprintf("DADD k%d k4 v4\r\n", i))
		if err != nil {
			panic(err.Error())
		}
	}
}
