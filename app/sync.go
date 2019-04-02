package app

import (
	"fmt"
	"math"
	"time"
)

const (
	// Approximate time (in seconds) to mine a block in mainnet
	MainNetTargetTimePerBlock = 300

	// Approximate time (in seconds) to mine a block in testnet
	TestNetTargetTimePerBlock = 120

	// Use 10% of estimated total headers fetch time to estimate rescan time
	RescanPercentage = 0.1

	// Use 80% of estimated total headers fetch time to estimate address discovery time
	DiscoveryPercentage = 0.8
)

type FetchHeadersData struct {
	StartHeaderHeight int32
	CurrentHeaderHeight int32
	BeginFetchTimeStamp int64
}

type FetchHeadersProgressReport struct {
	FetchedHeadersCount int32
	LastHeaderTime int64
	EstimatedBlocksToFetch int32
}

func UpdateFetchHeadersProgress(fetchHeadersData *FetchHeadersData, report FetchHeadersProgressReport, syncInfo *syncInfo) {
	// increment current block height value
	fetchHeadersData.CurrentHeaderHeight += report.FetchedHeadersCount

	// calculate percentage progress and eta
	totalFetchedHeaders := fetchHeadersData.CurrentHeaderHeight
	if fetchHeadersData.StartHeaderHeight > 0 {
		totalFetchedHeaders -= fetchHeadersData.StartHeaderHeight
	}

	syncEndPoint := report.EstimatedBlocksToFetch - fetchHeadersData.StartHeaderHeight
	headersFetchingRate := float64(totalFetchedHeaders) / float64(syncEndPoint)

	timeTakenSoFar := time.Now().Unix() - fetchHeadersData.BeginFetchTimeStamp
	estimatedTotalHeadersFetchTime := math.Round(float64(timeTakenSoFar) / headersFetchingRate)

	// 10% of estimated fetch time is used for estimating rescan time while 80% is used for estimating address discovery time
	estimatedRescanTime := estimatedTotalHeadersFetchTime * RescanPercentage
	estimatedDiscoveryTime := estimatedTotalHeadersFetchTime * DiscoveryPercentage
	estimatedTotalSyncTime := estimatedTotalHeadersFetchTime + estimatedRescanTime + estimatedDiscoveryTime

	totalTimeRemaining := (int64(estimatedTotalSyncTime) - timeTakenSoFar) / 60
	totalSyncProgress := (float64(timeTakenSoFar) / float64(estimatedTotalSyncTime)) * 100.0

	// update sync info
	syncInfo.FetchedHeadersCount = totalFetchedHeaders
	syncInfo.TotalHeadersToFetch = syncEndPoint
	syncInfo.HeadersFetchProgress = int32(math.Round(headersFetchingRate * 100))
	syncInfo.TotalTimeRemaining = fmt.Sprintf("%d min", totalTimeRemaining)
	syncInfo.TotalSyncProgress = int32(math.Round(totalSyncProgress))

	// calculate block header time difference
	hoursBehind := float64(time.Now().Unix() - report.LastHeaderTime) / 60
	daysBehind := int(math.Round(hoursBehind / 24))
	if daysBehind < 1 {
		syncInfo.DaysBehind = "<1 day"
	} else if daysBehind == 1 {
		syncInfo.DaysBehind = "1 day"
	} else {
		syncInfo.DaysBehind = fmt.Sprintf("%d days", daysBehind)
	}
}

func EstimateBlocksCount(netType string, bestBlockTimeStamp int64, bestBlock int32) int32 {
	var targetTimePerBlock int32
	if netType == "mainnet" {
		targetTimePerBlock = MainNetTargetTimePerBlock
	} else {
		targetTimePerBlock = TestNetTargetTimePerBlock
	}

	timeDifference := time.Now().Unix() - bestBlockTimeStamp
	return (int32(timeDifference) / targetTimePerBlock) + bestBlock
}
