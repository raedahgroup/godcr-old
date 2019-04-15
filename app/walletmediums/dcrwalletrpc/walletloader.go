package dcrwalletrpc

import (
	"context"
	"errors"
	"fmt"
	"github.com/decred/dcrd/dcrutil"
	"github.com/decred/dcrd/hdkeychain"
	"github.com/decred/dcrwallet/rpc/walletrpc"
	"github.com/decred/dcrwallet/walletseed"
	"github.com/raedahgroup/dcrlibwallet/blockchainsync"
	"github.com/raedahgroup/godcr/app/walletcore"
	"google.golang.org/grpc/codes"
)

var numberOfPeers int32

func (c *WalletRPCClient) GenerateNewWalletSeed() (string, error) {
	seed, err := hdkeychain.GenerateSeed(hdkeychain.RecommendedSeedLen)
	if err != nil {
		return "", err
	}

	return walletseed.EncodeMnemonic(seed), nil
}

func (c *WalletRPCClient) WalletExists() (bool, error) {
	res, err := c.walletLoader.WalletExists(context.Background(), &walletrpc.WalletExistsRequest{})
	if err != nil {
		return false, err
	}
	return res.Exists, nil
}

func (c *WalletRPCClient) CreateWallet(passphrase, seed string) error {
	seedBytes, err := walletseed.DecodeUserInput(seed)
	if err != nil {
		return err
	}

	_, err = c.walletLoader.CreateWallet(context.Background(), &walletrpc.CreateWalletRequest{
		PrivatePassphrase: []byte(passphrase),
		Seed:              seedBytes,
	})

	// wallet will be opened if the create operation was successful
	if err == nil {
		c.walletOpen = true
	}

	return err
}

func (c *WalletRPCClient) OpenWalletIfExist(ctx context.Context) (walletExists bool, err error) {
	c.walletOpen = false
	loadWalletDone := make(chan bool)

	go func() {
		defer func() {
			loadWalletDone <- true
		}()

		walletExists, err = c.WalletExists()
		if err != nil || !walletExists {
			return
		}

		_, err = c.walletLoader.OpenWallet(context.Background(), &walletrpc.OpenWalletRequest{})

		// ignore wallet already open errors, it could be that dcrwallet loaded the wallet when it was launched by the user
		// or godcr opened the wallet without closing it
		if isRpcErrorCode(err, codes.AlreadyExists) {
			err = nil
		}
	}()

	select {
	case <-loadWalletDone:
		// if err is nil, then wallet was opened
		c.walletOpen = err == nil
		return

	case <-ctx.Done():
		return false, ctx.Err()
	}
}

func (c *WalletRPCClient) IsWalletOpen() bool {
	// for now, assume that the wallet's already open since we're connecting through dcrwallet daemon
	// ideally, we'd have to use dcrwallet's WalletLoaderService to do this
	return c.walletOpen
}

func (c *WalletRPCClient) SyncBlockChain(showLog bool, syncInfoUpdated func(privateSyncInfo *blockchainsync.PrivateSyncInfo, updatedSection string)) error {
	ctx := context.Background()

	getBestBlock := func() int32 {
		bestBlockHeight, _ := c.BestBlock()
		return int32(bestBlockHeight)
	}
	getBestBlockTimestamp := func() int64 {
		var bestBlockInfo *walletrpc.BlockInfoResponse
		bestBlock, err := c.walletService.BestBlock(ctx, &walletrpc.BestBlockRequest{})
		if err == nil {
			bestBlockInfo, err = c.walletService.BlockInfo(ctx, &walletrpc.BlockInfoRequest{BlockHash: bestBlock.Hash})
		}
		if err != nil {
			return 0
		}
		return bestBlockInfo.Timestamp
	}

	// defaultSyncListener listens for reported sync updates, calculates progress and updates the caller via syncInfoUpdated
	defaultSyncListener := blockchainsync.DefaultSyncProgressListener(c.activeNet, showLog, getBestBlock, getBestBlockTimestamp, syncInfoUpdated)

	syncStream, err := c.walletLoader.SpvSync(ctx, &walletrpc.SpvSyncRequest{})
	if err != nil {
		return err
	}

	// read sync updates from syncStream in go routine and trigger defaultSyncListener methods to calculate progress and update caller
	go func() {
		for {
			syncUpdate, err := syncStream.Recv()
			if err != nil {
				defaultSyncListener.OnSyncError(blockchainsync.UnexpectedError, err)
				defaultSyncListener.OnSynced(false)
			}
			if syncUpdate.Synced {
				defaultSyncListener.OnSynced(true)
			}

			switch syncUpdate.NotificationType {
			case walletrpc.SyncNotificationType_FETCHED_HEADERS_STARTED:
				defaultSyncListener.OnFetchedHeaders(0, 0, blockchainsync.START)

			case walletrpc.SyncNotificationType_FETCHED_HEADERS_PROGRESS:
				defaultSyncListener.OnFetchedHeaders(syncUpdate.FetchHeaders.FetchedHeadersCount, syncUpdate.FetchHeaders.LastHeaderTime, blockchainsync.PROGRESS)

			case walletrpc.SyncNotificationType_FETCHED_HEADERS_FINISHED:
				defaultSyncListener.OnFetchedHeaders(0, 0, blockchainsync.FINISH)

			case walletrpc.SyncNotificationType_DISCOVER_ADDRESSES_STARTED:
				defaultSyncListener.OnDiscoveredAddresses(blockchainsync.START)

			case walletrpc.SyncNotificationType_DISCOVER_ADDRESSES_FINISHED:
				defaultSyncListener.OnDiscoveredAddresses(blockchainsync.FINISH)

			case walletrpc.SyncNotificationType_RESCAN_STARTED:
				defaultSyncListener.OnRescan(0, blockchainsync.START)

			case walletrpc.SyncNotificationType_RESCAN_PROGRESS:
				defaultSyncListener.OnRescan(syncUpdate.RescanProgress.RescannedThrough, blockchainsync.PROGRESS)

			case walletrpc.SyncNotificationType_RESCAN_FINISHED:
				defaultSyncListener.OnRescan(0, blockchainsync.FINISH)

			case walletrpc.SyncNotificationType_PEER_CONNECTED:
				numberOfPeers = syncUpdate.PeerInformation.PeerCount
				defaultSyncListener.OnPeerConnected(syncUpdate.PeerInformation.PeerCount)

			case walletrpc.SyncNotificationType_PEER_DISCONNECTED:
				numberOfPeers = syncUpdate.PeerInformation.PeerCount
				defaultSyncListener.OnPeerConnected(syncUpdate.PeerInformation.PeerCount)
			}
		}
	}()

	return nil
}

func (c *WalletRPCClient) RescanBlockChain() error {
	return nil // todo implement
}

func (c *WalletRPCClient) WalletConnectionInfo() (info walletcore.ConnectionInfo, err error) {
	accounts, loadAccountErr := c.AccountsOverview(walletcore.DefaultRequiredConfirmations)
	if loadAccountErr != nil {
		err = fmt.Errorf("error fetching account balance: %s", loadAccountErr.Error())
		info.TotalBalance = "0 DCR"
	} else {
		var totalBalance dcrutil.Amount
		for _, acc := range accounts {
			totalBalance += acc.Balance.Total
		}
		info.TotalBalance = totalBalance.String()
	}

	bestBlock, bestBlockErr := c.BestBlock()
	if bestBlockErr != nil && err != nil {
		err = fmt.Errorf("%s, error in fetching best block %s", err.Error(), bestBlockErr.Error())
	} else if bestBlockErr != nil {
		err = bestBlockErr
	}

	info.LatestBlock = bestBlock
	info.NetworkType = c.NetType()
	info.PeersConnected = numberOfPeers

	return
}

func (c *WalletRPCClient) BestBlock() (uint32, error) {
	req, err := c.walletService.BestBlock(context.Background(), &walletrpc.BestBlockRequest{})
	if err != nil {
		return 0, err
	}
	return req.Height, err
}

func (c *WalletRPCClient) CloseWallet() {
	// don't actually close wallet loaded by dcrwallet
	// - if wallet wasn't opened by godcr, closing it could cause troubles for user
	// - even if wallet was opened by godcr, closing it without closing dcrwallet would cause troubles for user
	// when they next launch godcr
}

func (c *WalletRPCClient) DeleteWallet() error {
	return errors.New("wallet cannot be deleted when connecting via dcrwallet rpc")
}
