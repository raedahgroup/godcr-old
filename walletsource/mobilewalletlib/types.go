package mobilewalletlib

type BlockChainSyncListener struct {
	SyncStarted         func()
	SyncEnded           func(err error)
	OnHeadersFetched    func(percentageProgress int64)
	OnDiscoveredAddress func(state string)
}
