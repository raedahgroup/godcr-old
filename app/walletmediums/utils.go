package walletmediums

import "time"

const (
	MainNetTargetTimePerBlock = 300
	TestNetTargetTimePerBlock = 120
)

func CalculateBlockSyncProgress(netType string, bestBlock, lastHeaderTime int64) int64 {
	var targetTimePerBlock int64
	if netType == "mainnet" {
		targetTimePerBlock = MainNetTargetTimePerBlock
	} else {
		targetTimePerBlock = TestNetTargetTimePerBlock
	}

	estimatedBlocks := ((time.Now().Unix() - lastHeaderTime) / targetTimePerBlock) + bestBlock
	fetchedPercentage := bestBlock / estimatedBlocks * 100

	if fetchedPercentage >= 100 {
		fetchedPercentage = 100
	}

	return fetchedPercentage
}
