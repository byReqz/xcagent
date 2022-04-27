package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"time"
)

var Back bool

func forkbg() {
	f, err := os.CreateTemp("/tmp", "xcagent-*.tmp")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	self, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(self, "-background", f.Name())
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	pid := cmd.Process.Pid
	fmt.Println(pid)
	err = cmd.Process.Release()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("XCAGENT_PID=" + fmt.Sprint(pid))
	err = os.Setenv("XCAGENT_PID", fmt.Sprint(pid))
	if err != nil {
		log.Fatal(err)
	}
}

func inbg() {
	time.Sleep(30 * time.Second)
	if len(flag.Args()) > 0 {
		f, err := os.ReadFile(flag.Args()[0])
		if err != nil {
			os.Exit(1)
		}
		err = os.Remove(flag.Args()[0])
		if err != nil {
			os.Exit(1)
		}
		//missing
	}
}

func init() {
	flag.BoolVar(&Back, "background", false, "sets background var")
	flag.Parse()
}

func main() {
	if !Back {
		fmt.Println("forking into bg")
		forkbg()
	} else if Back {
		inbg()
	} else {
		fmt.Println("exited")
	}
}
