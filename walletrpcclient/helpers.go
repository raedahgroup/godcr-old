package walletrpcclient

import (
	"github.com/decred/dcrd/dcrutil"
)

func AtomToCoin(atom int64) int64 {
	val := atom / dcrutil.AtomsPerCoin
	return int64(val)
}
