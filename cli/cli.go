package cli

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
	qrcode "github.com/skip2/go-qrcode"
)

type (
	response struct {
		columns []string
		result  [][]interface{}
	}
	// Handler carries out the action required by a command.
	// commandArgs holds the arguments passed to the command.
	Handler func(commandArgs []string) (*response, error)

	// CLI holds data needed to run the program.
	CLI struct {
		funcMap          map[string]Handler
		commands         map[string]string
		descriptions     map[string]string
		commandListOrder []string
		appName          string
		walletrpcclient  *walletrpcclient.Client
	}
)

// New creates a new CLI object with the given arguments.
func New(walletrpcclient *walletrpcclient.Client, appName string) *CLI {
	client := &CLI{
		funcMap:         make(map[string]Handler),
		commands:        make(map[string]string),
		descriptions:    make(map[string]string),
		walletrpcclient: walletrpcclient,
		appName:         appName,
	}

	// register handlers
	client.registerHandlers()

	return client
}

func (c *CLI) registerHandlers() {
	c.registerHandler("balance", "balance", "show your balance", c.balance)
	c.registerHandler("send", "send", "send a transaction", c.send)
	c.registerHandler("receive", "receive", "show your address to receive funds", c.receive)
}

// registerHandler registers a command, its description and its handler
func (c *CLI) registerHandler(key, command, description string, h Handler) {
	if _, ok := c.funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	c.funcMap[key] = h
	c.commands[key] = command
	c.descriptions[key] = description
	c.commandListOrder = append(c.commandListOrder, key)
}

func (c *CLI) balance(commandArgs []string) (*response, error) {
	balances, err := c.walletrpcclient.Balance()
	if err != nil {
		return nil, err
	}

	res := &response{
		columns: []string{
			"Account",
			"Total",
			"Spendable",
			"Locked By Tickets",
			"Voting Authority",
			"Unconfirmed",
		},
		result: make([][]interface{}, len(balances)),
	}
	for i, v := range balances {
		res.result[i] = []interface{}{
			v.AccountName,
			v.Total,
			v.Spendable,
			v.LockedByTickets,
			v.VotingAuthority,
			v.Unconfirmed,
		}
	}

	return res, nil
}

func (c *CLI) send(commandArgs []string) (*response, error) {
	sourceAccount, err := getSendSourceAccount(c.walletrpcclient)
	if err != nil {
		return nil, err
	}

	destinationAddress, err := getSendDestinationAddress(c.walletrpcclient)
	if err != nil {
		return nil, err
	}

	sendAmount, err := getSendAmount()
	if err != nil {
		return nil, err
	}

	passphrase, err := getWalletPassphrase()
	if err != nil {
		return nil, err
	}

	result, err := c.walletrpcclient.Send(sendAmount, sourceAccount, destinationAddress, passphrase)
	if err != nil {
		return nil, err
	}

	res := &response{
		columns: []string{
			"Result",
			"Hash",
		},
		result: [][]interface{}{
			[]interface{}{
				"The transaction was published successfully",
				result.TransactionHash,
			},
		},
	}

	return res, nil
}

func (c *CLI) receive(commandArgs []string) (*response, error) {
	var recieveAddress uint32 = 0

	// if no address passed in
	if len(commandArgs) == 0 {

		// display menu options to select account
		var err error
		recieveAddress, err = getSendSourceAccount(c.walletrpcclient)
		if err != nil {
			return nil, err
		}
	} else {
		// if an address was passed in eg. ./dcrcli receive 0 use that address
		x, err := strconv.ParseUint(commandArgs[0], 10, 32)
		if err != nil {
			return nil, fmt.Errorf("Error parsing account number: %s", err.Error())
		}

		recieveAddress = uint32(x)
	}

	r, err := c.walletrpcclient.Receive(recieveAddress)
	if err != nil {
		return nil, err
	}

	qr, err := qrcode.New(r.Address, qrcode.Medium)
	if err != nil {
		return nil, fmt.Errorf("Error generating QR Code: %s", err.Error())
	}

	res := &response{
		columns: []string{
			"Address",
			"QR Code",
		},
		result: [][]interface{}{
			[]interface{}{
				r.Address,
				qr.ToString(true),
			},
		},
	}
	return res, nil
}

// RunCommand invokes the handler function registered for the given
// command in `commandArgs`.
//
// If no command, or an unsupported command is passed to RunCommand,
// the program exits with an error.
// commandArgs[0] is the command to run. commandArgs[1:] are the arguments to the command.
func (c *CLI) RunCommand(commandArgs []string) {
	if len(commandArgs) == 0 {
		c.noCommandReceived()
		os.Exit(1)
	}

	command := commandArgs[0]
	if !c.isCommandSupported(command) {
		c.invalidCommandReceived(command)
		os.Exit(1)
	}

	handler := c.funcMap[command]
	res, err := handler(commandArgs[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command '%s'\n", command)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	printResult(res)
	os.Exit(0)
}

func (c *CLI) noCommandReceived() {
	fmt.Printf("usage: %s [OPTIONS] <command> [<args...>]\n\n", c.appName)
	fmt.Printf("available %s commands:\n", c.appName)
	c.listCommands()
}

func (c *CLI) listCommands() {
	res := &response{
		columns: []string{"Command", "Description"},
	}

	for _, key := range c.commandListOrder {
		item := []interface{}{
			c.commands[key],
			c.descriptions[key],
		}

		res.result = append(res.result, item)
	}
	printResult(res)
}

// IsCommandSupported returns true if the `command` specified is registered
// on the current CLI object; otherwise, it returns false.
func (c *CLI) isCommandSupported(command string) bool {
	_, ok := c.funcMap[command]
	return ok
}

func (c *CLI) invalidCommandReceived(command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a supported command.\n\n", c.appName, command)
	c.noCommandReceived()
}

func printResult(res *response) {
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
	header := ""
	spaceRow := ""
	columnLength := len(res.columns)

	for i := range res.columns {
		tab := " \t "
		if columnLength == i+1 {
			tab = " "
		}
		header += res.columns[i] + tab
		spaceRow += " " + tab
	}

	fmt.Fprintln(w, header)
	fmt.Fprintln(w, spaceRow)
	for _, row := range res.result {
		rowStr := ""
		for range row {
			rowStr += "%v \t "
		}

		rowStr = strings.TrimSuffix(rowStr, "\t ")
		fmt.Fprintln(w, fmt.Sprintf(rowStr, row...))
	}

	w.Flush()
}
