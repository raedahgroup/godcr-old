package cli

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"text/tabwriter"

	"github.com/raedahgroup/dcrcli/walletrpcclient"
)

type (
	response struct {
		columns []string
		result  [][]interface{}
	}
	Handler func(commandArgs []string) (*response, error)
	CLI     struct {
		funcMap         map[string]Handler
		commands        map[string]string
		descriptions    map[string]string
		appName         string
		walletrpcclient *walletrpcclient.Client
	}
)

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

func (c *CLI) RunCommand(commandArgs []string) {
	if len(commandArgs) == 0 {
		c.noCommandReceived()
		os.Exit(1)
	}

	command := commandArgs[0]
	if command == "-l" {
		command = "listcommands"
	} else if !c.IsCommandSupported(command) {
		c.invalidCommandReceived(command)
		os.Exit(1)
	}

	handler := c.funcMap[command]
	res, err := handler(commandArgs[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error executing command '%s'", command)
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	printResult(res)
	os.Exit(0)
}

func (c *CLI) noCommandReceived() {
	fmt.Printf("usage: %s [OPTIONS] <command> [<args...>]\n\n", c.appName)
	fmt.Printf("available %s commands:\n", c.appName)
	c.listCommands(nil)
	fmt.Printf("\nFor available options, see '%s -h'\n", c.appName)
}

func (c *CLI) invalidCommandReceived(command string) {
	fmt.Fprintf(os.Stderr, "%s: '%s' is not a valid command. See '%s -h'\n", c.appName, command, c.appName)
}

func (c *CLI) IsCommandSupported(command string) bool {
	_, ok := c.funcMap[command]
	return ok
}

// RegisterHandler registers a command, its description and its handler
func (c *CLI) RegisterHandler(key, command, description string, h Handler) {
	if _, ok := c.funcMap[key]; ok {
		panic("trying to register a handler twice: " + key)
	}

	c.funcMap[key] = h
	c.commands[key] = command
	c.descriptions[key] = description
}

func (c *CLI) registerHandlers() {
	c.RegisterHandler("listcommands", "-l", "List all supported commands", c.listCommands)
	c.RegisterHandler("receive", "receive", "Generate address to receive funds", c.receive)
	c.RegisterHandler("send", "send", "Send DCR to address. Multi-step", c.send)
	c.RegisterHandler("balance", "balance", "Check balance of an account", c.balance)
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

func (c *CLI) listCommands(commandArgs []string) (*response, error) {
	res := &response{
		columns: []string{"Command", "Description"},
	}

	for i, v := range c.commands {
		item := []interface{}{
			v,
			c.descriptions[i],
		}

		res.result = append(res.result, item)
	}
	return res, nil
}

func (c *CLI) receive(commandArgs []string) (*response, error) {
	if len(commandArgs) == 0 {
		return nil, errors.New(" account number is required.  \nUsage:\n  receive \"accountnumber\"")
	}

	acc, err := strconv.ParseUint(commandArgs[0], 0, 32)
	if err != nil {
		return nil, fmt.Errorf("Error parsing account number: %s", err.Error())
	}

	r, err := c.walletrpcclient.Receive(uint32(acc))
	if err != nil {
		return nil, err
	}

	res := &response{
		columns: []string{
			"Address",
		},
		result: [][]interface{}{
			[]interface{}{
				r.Address,
			},
		},
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
