// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package nuklog

import "fmt"

func LogInfo(message string) {
	log.Info(message)
	fmt.Println(message)
}

func LogWarn(message string) {
	log.Warn(message)
	fmt.Println(message)
}

func LogError(message error) {
	log.Error(message)
	fmt.Println(message)
}
