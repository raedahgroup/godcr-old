package routes

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/raedahgroup/godcr/app/walletcore"
	"github.com/raedahgroup/godcr/web/weblog"
)

var clients = make(map[*websocket.Conn]bool)
var wsBroadcast = make(chan Packet)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type eventType string

const (
	UpdateConnectionInfo eventType = "updateConnInfo"
	UpdateBalance        eventType = "updateBalance"
)

type Packet struct {
	Event   eventType   `json:"event"`
	Message interface{} `json:"message"`
}

func (routes *Routes) wsHandler(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	clients[ws] = true
}

func waitToSendMessagesToClients() {
	for {
		msg := <-wsBroadcast
		for client := range clients {
			err := client.WriteJSON(msg)
			if err != nil {
				log.Printf("error: %s", err.Error())
				client.Close()
				delete(clients, client)
			}
		}
	}
}

func (routes *Routes) sendWsConnectionInfoUpdate() {
	var info = walletcore.ConnectionInfo{
		NetworkType:    routes.walletMiddleware.NetType(),
		PeersConnected: routes.walletMiddleware.GetConnectedPeersCount(),
	}

	wsBroadcast <- Packet{
		Event:   UpdateConnectionInfo,
		Message: info,
	}
}

func (routes *Routes) sendWsBalance() {
	accounts, err := routes.walletMiddleware.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if err != nil {
		weblog.LogError(fmt.Errorf("Error fetching account balance: %s", err.Error()))
		return
	}
	type accountInfo struct {
		Number uint32 `json:"number"`
		Balance string `json:"balance"`
	}

	var accountInfos []accountInfo

	var totalBalance walletcore.Balance
	for _, acc := range accounts {
		totalBalance.Spendable += acc.Balance.Spendable
		totalBalance.Total += acc.Balance.Total

		accountInfos = append(accountInfos, accountInfo{Number: acc.Number, Balance: acc.Balance.Total.String()})
	}
	wsBroadcast <- Packet{
		Event:   UpdateBalance,
		Message: map[string]interface{}{"accounts": accountInfos, "total": totalBalance.String()},
	}
}
