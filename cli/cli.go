package cli

import (
	"bytes"
	"fmt"
	"io"
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
	// handler carries out the action required by a command.
	// commandArgs holds the arguments passed to the command.
	handler func(walletrpcclient *walletrpcclient.Client, commandArgs []string) (*response, error)

	// cli holds data needed to run the program.
	cli struct {
		funcMap          map[string]handler
		commandListOrder []string
		appName          string
		walletrpcclient  *walletrpcclient.Client
	}

	// command is an action that can be requested at the cli
	command struct {
		name        string
		description string
		handler     handler
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

// supportedCommands provides the commands available from the cli
func supportedCommands() []command {
	return []command{
		{"balance", "show your balance", balance},
		{"send", "send a transaction", send},
		{"receive", "show your address to receive funds", receive},
	}
}

func (c *cli) registerHandlers() {
	commands := supportedCommands()
	for _, command := range commands {
		c.registerHandler(command.name, command.name, command.description, command.handler)
	}
}

// registerHandler registers a command, its description and its handler
func (c *cli) registerHandler(key, command, description string, h handler) {
	if _, ok := c.funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	c.funcMap[key] = h
	c.commandListOrder = append(c.commandListOrder, key)
}

func balance(walletrpcclient *walletrpcclient.Client, commandArgs []string) (*response, error) {
	balances, err := walletrpcclient.Balance()
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

func send(walletrpcclient *walletrpcclient.Client, commandArgs []string) (*response, error) {
	sourceAccount, err := getSendSourceAccount(walletrpcclient)
	if err != nil {
		return nil, err
	}

	destinationAddress, err := getSendDestinationAddress(walletrpcclient)
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

	result, err := walletrpcclient.SendFromAccount(sendAmount, sourceAccount, destinationAddress, passphrase)
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

func receive(walletrpcclient *walletrpcclient.Client, commandArgs []string) (*response, error) {
	var recieveAddress uint32 = 0

	// if no address passed in
	if len(commandArgs) == 0 {

		// display menu options to select account
		var err error
		recieveAddress, err = getSendSourceAccount(walletrpcclient)
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

	r, err := walletrpcclient.Receive(recieveAddress)
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

// UsageString returns the cli usage message as a string.
func UsageString() string {
	buf := bytes.NewBuffer([]byte{})
	recordHelpMessage(buf)
	return strings.TrimPrefix(buf.String(), usagePrefix)
}

// recordHelpMessage outputs help message to a buffer
func recordHelpMessage(buf io.Writer) {
	outputDest := helpTabWriter(buf)
	writeHelpMessage(outputDest)
}

var (
	stderrHelpWriter = helpTabWriter(os.Stderr)
	usagePrefix      = "Usage:\n  dcrcli "
)

func helpTabWriter(w io.Writer) *tabwriter.Writer {
	return tabwriter.NewWriter(w, 0, 0, 1, ' ', tabwriter.AlignRight|tabwriter.Debug)
}

// PrintHelp outputs help message to os.Stderr
func PrintHelp() {
	writeHelpMessage(stderrHelpWriter)
}

func writeHelpMessage(w *tabwriter.Writer) {
	res := &response{
		columns: []string{usagePrefix + "[OPTIONS] <command> [<args...>]\n\nAvailable commands:"},
	}
	commands := supportedCommands()

	for _, command := range commands {
		item := []interface{}{
			command.name,
			command.description,
		}
		res.result = append(res.result, item)
	}

	printResult(w, res)
}

// RunCommand invokes the handler function registered for the given
// command in `commandArgs`.
//
// If no command, or an unsupported command is passed to RunCommand,
// the program exits with an error.
// commandArgs[0] is the command to run. commandArgs[1:] are the arguments to the command.
func (c *cli) RunCommand(commandArgs []string) {
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
	res, err := handler(c.walletrpcclient, commandArgs[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command '%s'\n", command)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	printResult(stderrHelpWriter, res)
	os.Exit(0)
}

func (c *cli) noCommandReceived() {
	PrintHelp()
}

// IsCommandSupported returns true if the `command` specified is registered
// on the current cli object; otherwise, it returns false.
func (c *cli) isCommandSupported(command string) bool {
	_, ok := c.funcMap[command]
	return ok
}

func (c *cli) invalidCommandReceived(command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a supported command.\n\n", c.appName, command)
	PrintHelp()
}

func printResult(w *tabwriter.Writer, res *response) {
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
