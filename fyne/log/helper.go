// Copyright (c) 2013-2014 The btcsuite developers
// Copyright (c) 2018 The Decred developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package log

import "fmt"

func Info(messages ...string) {
	for _, message := range messages {
		logger.Info(message)
	}
}

func Warn(messages ...string) {
	for _, message := range messages {
		logger.Warn(message)
	}
}

func Error(messages ...string) {
	for _, message := range messages {
		logger.Error(message)
	}
}

func PrintInfo(messages ...string) {
	Info(messages...)
	fmt.Println(messages)
}

func PrintWarn(messages ...string) {
	Warn(messages...)
	fmt.Println(messages)
}

func PrintError(messages ...string) {
	Error(messages...)
	fmt.Println(messages)
}
