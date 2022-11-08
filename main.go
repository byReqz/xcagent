package main

import (
	"github.com/byReqz/xcagent/client"
	"github.com/byReqz/xcagent/daemon"
	flag "github.com/spf13/pflag"
	"os"
)

var (
	Daemon bool
)

func init() {
	var quiet bool

	flag.BoolVarP(&quiet, "quiet", "q", false, "disables output")
	flag.BoolVar(&Daemon, "daemon", false, "starts xcagent in daemon mode")
	flag.Parse()

	if quiet {
		os.Stdout = nil
		os.Stderr = nil
	}
}

func main() {
	if Daemon {
		daemon.Listener()
	} else {
		client.Main()
	}
}
