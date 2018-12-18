package cli

import (
	"fmt"
	"os"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type (
	response struct {
		columns []string
		result  [][]interface{}
	}
	// handler carries out the action required by a command.
	// commandArgs holds the arguments passed to the command.
	handler func(c *cli, commandArgs []string) (*response, error)

	// cli holds data needed to run the program.
	cli struct {
		funcMap         map[string]handler
		appName         string
		walletrpcclient *walletrpcclient.Client
	}
)

// New creates a new cli object with the given arguments.
func New(walletrpcclient *walletrpcclient.Client, appName string) *cli {
	client := &cli{
		funcMap:         make(map[string]handler),
		walletrpcclient: walletrpcclient,
		appName:         appName,
	}

	// register handlers
	client.registerHandlers()

	return client
}

func (c *cli) registerHandlers() {
	commands := supportedCommands()
	for _, command := range commands {
		c.registerHandler(command.name, command.handler)
	}
}

// registerHandler registers a command, its description and its handler
func (c *cli) registerHandler(key string, h handler) {
	if _, ok := c.funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	c.funcMap[key] = h
}

// RunCommand invokes the handler function registered for the given
// command in `commandArgs`.
//
// If no command, or an unsupported command is passed to RunCommand,
// the program exits with an error.
// commandArgs[0] is the command to run. commandArgs[1:] are the arguments to the command.
func (c *cli) RunCommand(commandArgs []string) {
	if len(commandArgs) == 0 {
		writeSimpleHelpMessage()
		os.Exit(1)
	}

	command := commandArgs[0]
	if !c.isCommandSupported(command) {
		c.invalidCommandReceived(command)
		os.Exit(1)
	}

	handler := c.funcMap[command]
	res, err := handler(c, commandArgs[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command '%s'\n", command)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	printResult(tabWriter(os.Stdout), res)
	os.Exit(0)
}

// IsCommandSupported returns true if the `command` specified is registered
// on the current cli object; otherwise, it returns false.
func (c *cli) isCommandSupported(command string) bool {
	_, ok := c.funcMap[command]
	return ok
}

func (c *cli) invalidCommandReceived(command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a supported command.\n", c.appName, command)
}
