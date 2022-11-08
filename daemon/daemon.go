package daemon

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
)

var (
	Passphrase string // passphrase of the DB
	Path       string // path to the DB
)

// callKeepass calls keepassxc-cli with the given password and args. Returns stdout or stderr.
func callKeepass(args ...string) (string, error) {
	cmd := exec.Command("keepassxc-cli", args...)
	cmd.Stdin = strings.NewReader(Passphrase)
	out, err := cmd.Output()
	if err == nil {
		return string(out), nil
	}
	exitErr, isExitErr := err.(*exec.ExitError)
	if isExitErr {
		return string(exitErr.Stderr), err
	}
	return "", err
}

// ReadConn reads the connection till the defined delimiter (0x04 / EOT).
func ReadConn(c net.Conn) (string, error) {
	var buf string
	for {
		r := make([]byte, 1)
		_, err := c.Read(r)
		if err != nil {
			return "", err
		}
		if r[0] == 4 { // EOT as query delimiter
			break
		} else {
			buf = buf + string(r)
		}
	}
	return buf, nil
}

// WriteConn sends data into the connection and terminates it by sending the given delimiter (0x04 / EOT).
func WriteConn(c net.Conn, data string) error {
	_, err := c.Write([]byte(data))
	if err != nil {
		return err
	}
	_, err = c.Write([]byte{4})
	return err
}

// command executes the given arguments for keepassxc-cli and returns either stdout, stderr or the actual error.
func command(args ...string) string {
	out, err := callKeepass(args...)
	if out != "" {
		return out
	}
	if err != nil {
		return err.Error()
	}
	return ""
}

// setPassphrase saves the supplied credentials into memory.
func setPassphrase(path, key string) {
	Path = path
	Passphrase = key
}

// connHandler handles the individual connections coming into the socket.
func connHandler(l net.Listener, c net.Conn) {
	defer func() { _ = c.Close() }()
	buf, err := ReadConn(c)
	if err != nil {
		_ = WriteConn(c, "ERROR: failed reading request") // best effort to try to inform the other side of the issue
		return
	}

	args := strings.Split(buf, " ")
	if len(args) == 0 {
		_ = WriteConn(c, "ERROR: empty request") // best effort to try to inform the other side of the issue
		return
	}

	switch args[0] {
	case "KILLAGENT":
		_ = WriteConn(c, "OK")
		_ = c.Close()
		_ = l.Close()
		os.Exit(0)
	case "PING":
		_ = WriteConn(c, "PONG")
	case "COMMAND":
		if len(args) < 2 {
			_ = WriteConn(c, "ERROR: no command given")
		} else {
			_ = WriteConn(c, command(args[1:]...))
		}
	case "SET-PASSPHRASE":
		if len(args) < 3 {
			_ = WriteConn(c, "ERROR: missing arguments")
		} else {
			setPassphrase(args[1], args[2])
			_ = WriteConn(c, "OK")
		}
	case "QUERY":
		if len(args) < 2 {
			_ = WriteConn(c, "ERROR: missing argument")
		} else {
			_ = WriteConn(c, command("show", "-q", "-s", "--all", Path, args[1]))
		}
	default:
		_ = WriteConn(c, "ERROR: unknown command: "+args[0])
	}
	_ = c.Close()
}

// Listener is the actual daemon listening on the socket.
func Listener() {
	l, err := net.Listen("unix", "/tmp/xcagent-"+fmt.Sprint(os.Getpid())+".sock")
	if err != nil {
		log.Fatal(err)
	}
	defer func() { _ = l.Close() }()

	fmt.Println(l)

	for {
		c, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}
		go connHandler(l, c)
	}
}

// Init starts the daemon in the background and detaches it from the current process. Will print environment variables for the user to set.
func Init() {
	self, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}
	cmd := exec.Command(self, "--daemon")
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	pid := cmd.Process.Pid
	err = cmd.Process.Release()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("XCAGENT_SOCK=/tmp/xcagent-%v.sock; export XCAGENT_SOCK;\n", pid)
	fmt.Printf("XCAGENT_PID=%v; export XCAGENT_PID;\n", pid)
	fmt.Printf("echo Agent pid %v;\n", pid)
}
