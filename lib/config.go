// Package lib provides ...
package lib

import (
	"flag"
)

type cliParams struct {
	ServerAddr string
	LogEnable bool
	LogLevel int
	LogPath string
	CpuProf string
	ProfDir string
	ExitOn int
}

var CliParams *cliParams = &cliParams{}

func init() {
	flag.StringVar(&CliParams.ServerAddr, "addr", "127.0.0.1:8800", "Address to use by server")
	flag.BoolVar(&CliParams.LogEnable, "log", false, "Log on/off")
	flag.IntVar(&CliParams.LogLevel, "log-level", 1, "Log level [1-5]")
	flag.StringVar(&CliParams.LogPath, "log-path", "", "Path to logs dir")
	flag.StringVar(&CliParams.CpuProf, "cpu-prof", "", "Path to cpu.pprof file")
	flag.StringVar(&CliParams.CpuProf, "prof-dir", "", "Path to profile directory")
	flag.IntVar(&CliParams.ExitOn, "exit-on", 0, "Automatically stop app after N sec")
	flag.Parse()
}

