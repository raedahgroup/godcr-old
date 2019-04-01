package app

import "time"

const (
	MainNetTargetTimePerBlock = 300
	TestNetTargetTimePerBlock = 120
)

func EstimateBlocksCount(netType string, bestBlockTimeStamp, lastHeaderTimeStamp int64) int64 {
	var targetTimePerBlock int64
	if netType == "mainnet" {
		targetTimePerBlock = MainNetTargetTimePerBlock
	} else {
		targetTimePerBlock = TestNetTargetTimePerBlock
	}

	return ((time.Now().Unix() - lastHeaderTimeStamp) / targetTimePerBlock) + bestBlockTimeStamp
}
