package app

import "time"

const (
	MainNetTargetTimePerBlock = 300
	TestNetTargetTimePerBlock = 120
)

func EstimateBlocksCount(netType string, bestBlockTimeStamp int64, bestBlock int32) int64 {
	var targetTimePerBlock int64
	if netType == "mainnet" {
		targetTimePerBlock = MainNetTargetTimePerBlock
	} else {
		targetTimePerBlock = TestNetTargetTimePerBlock
	}

	return ((time.Now().Unix() - bestBlockTimeStamp) / targetTimePerBlock) + int64(bestBlock)
}
