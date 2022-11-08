package client

import (
	"fmt"
	gap "github.com/byReqz/go-ask-password"
	"github.com/byReqz/xcagent/daemon"
	flag "github.com/spf13/pflag"
	"log"
	"net"
	"os"
	"strings"
)

func printHelp() {
	help := `xcagent: help

arguments:
	start                       -    start the daemon in the background
	kill                        -    kill the background agent
	ping                        -    ping the background agent
	command <args...>           -    send a raw command to keepassxc-cli through the agent
	query <entry name>          -    query for an entry through the agent
	set-passphrase <path to db> -    set database path and passphrase for the agent

flags:
	--daemon        -    start daemon mode in foreground
	-q / --quiet    -    disable output`
	fmt.Println(help)
}

// sendSocket sends the given query to the socket address thats set as XCAGENT_SOCK environment variable and returns the reply.
func sendSocket(args ...string) {
	addr := os.Getenv("XCAGENT_SOCK")
	if addr == "" {
		log.Fatal("could not read socket address from environment")
	}
	conn, err := net.Dial("unix", addr)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(strings.Join(args, " ")))
	if err != nil {
		log.Fatal(err)
	}
	conn.Write([]byte{4})
	reply, err := daemon.ReadConn(conn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(strings.TrimSpace(reply))
}

// Main handles the arguments passed to the client executable.
func Main() {
	args := flag.Args()
	if len(args) == 0 {
		printHelp()
		os.Exit(0)
	}

	switch args[0] {
	case "start":
		daemon.Init()
	case "kill":
		sendSocket("KILLAGENT")
	case "ping":
		sendSocket("PING")
	case "command":
		if len(args) < 2 {
			log.Fatal("no command given")
		}
		sendSocket("COMMAND", args[1])
	case "set-passphrase":
		if len(args) < 2 {
			log.Fatal("no database given")
		}
		key, err := gap.AskPassword("Passphrase: ")
		if err != nil {
			log.Fatal(err)
		}
		sendSocket("SET-PASSPHRASE", args[1], key)
	case "query":
		if len(args) < 2 {
			log.Fatal("no query given")
		}
		sendSocket("QUERY", args[1])
	default:
		fmt.Println("[Unknown Command]")
		printHelp()
	}
}
