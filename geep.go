// Package main provides profiling tools
package main

import (
	"os"
	"fmt"
	"time"
	"runtime"
	"runtime/pprof"
	server "github.com/avsolo/gache/server"
	ll "github.com/avsolo/gache/lib"
	// "github.com/pkg/profile"
)

func init() {
	c := runtime.NumCPU()
    runtime.GOMAXPROCS(c)
	fmt.Printf("Using CPU num: %d\n", c)
}

func main() {
	if ll.CliParams.CpuProf != "" {
		f, err := os.Create(ll.CliParams.CpuProf)
		if err != nil {
			fmt.Printf("Can't open profiler file. Error: %s\n", err.Error())
			return
		}
		pprof.StartCPUProfile(f)
		fmt.Printf("Starting profile to: %s\n", f.Name())
		defer func() {
			fmt.Printf("Stop profiler\n")
			pprof.StopCPUProfile()
		}()

		// Not working on OS X
		// defer profile.Start(
			// profile.CPUProfile,
			// profile.ProfilePath(ll.CliParams.ProfDir)).Stop()

	} else {
		fmt.Printf("Profile disabled\n")
	}

	// Sample func for stopping app after N sec
	if ll.CliParams.ExitOn > 0 {
		fmt.Printf("Exit after %d sec\n", ll.CliParams.ExitOn)
		go func() {
			t := time.Now().Unix() + int64(ll.CliParams.ExitOn)
			for {
				c := time.Now().Unix()
				if c >= t {
					fmt.Printf("Exit\n")
					os.Exit(0)
				}
				time.Sleep(time.Duration(1) * time.Second)
			}
		}()
	} else {
		fmt.Printf("Running forever\n")
	}

	srv := server.NewServer(ll.CliParams.ServerAddr)
	srv.ListenTCP()
}
