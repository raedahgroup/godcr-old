package walletrpcclient

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func sendTransaction(conn *grpc.ClientConn, ctx context.Context, opts []string) (*Response, error) {
	//c := pb.NewWalletServiceClient(conn)
	var fromAccount int32
	fromAccountFunc := func() error {
		fmt.Println("From Account: ")
		_, err := fmt.Scanf("%d", &fromAccount)
		if err != nil {
			return fmt.Errorf("Error reading input: %s", err.Error())
		}
		return nil
	}

	for {
		err := fromAccountFunc()
		if err == nil {
			break
		}
		fmt.Println(err.Error())
	}

	var toAddress string
	toAddressFunc := func() error {
		fmt.Println("Destination Address: ")
		_, err := fmt.Scanln(&toAddress)
		if err != nil {
			return fmt.Errorf("Error reading input: %s", err.Error())
		}
		return nil
	}

	for {
		err := toAddressFunc()
		if err == nil {
			break
		}
		fmt.Println(err.Error())
	}

	var amount int64
	amountFunc := func() error {
		fmt.Println("Send Amount: ")
		_, err := fmt.Scanf("%d", &amount)
		if err != nil {
			return fmt.Errorf("Error reading input: %s", err.Error())
		}
		return nil
	}

	for {
		err := amountFunc()
		if err != nil {
			break
		}
		fmt.Println(err.Error())
	}

	return nil, nil
}

/**
reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter text: ")
    text, _ := reader.ReadString('\n')
    fmt.Println(text)

    fmt.Println("Enter text: ")
    text2 := ""
    fmt.Scanln(text2)
    fmt.Println(text2)

    ln := ""
    fmt.Sscanln("%v", ln)
    fmt.Println(ln)
func getTransactions(conn *grpc.ClientConn, ctx context.Context, opts []string) (*Response, error) {
	c := pb.NewWalletServiceClient(conn)

	// check if passed options are complete
	if len(opts) < 2 {
		return nil, fmt.Errorf("command 'transactions' requires 2 params. %d found", len(opts))
	}

	// get block height
	startingBlockheight, err := strconv.ParseInt(opts[0], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error getting starting block height from options: %s", err.Error())
	}

	limit, err := strconv.ParseInt(opts[1], 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error getting limit from options: %s", err.Error())
	}

	req := &pb.GetTransactionsrequest{}
}
**/
