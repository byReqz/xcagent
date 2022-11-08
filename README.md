# xcagent
password retention agent for keepassxc

## Building
xcagent can be built like any other go app:
```
CGO_ENABLED=0 go build -ldflags="-s -w" -o xcagent ./main.go 
```

## Usage
The usage is similar to `ssh-agent`, in that `xcagent start` prints environment values that are used to query the agent later on. The default flow to set up xcagent is therefore:
1. `eval $(xcagent start)`
2. `xcagent set-passphrase <path to database>`
3. `xcagent query <entry name>` to get the data for the wished entry

```
xcagent: help

arguments:
start                       -    start the daemon in the background
kill                        -    kill the background agent
ping                        -    ping the background agent
command <args...>           -    send a raw command to keepassxc-cli through the agent
query <entry name>          -    query for an entry through the agent
set-passphrase <path to db> -    set database path and passphrase for the agent

flags:
--daemon        -    start daemon mode in foreground
-q / --quiet    -    disable output
```