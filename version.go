package main

import (
	"fmt"
	"runtime/debug"
	"strings"
	"sync"
)

var (
	sdkVersion     = "8.unknown"
	sdkVersionOnce sync.Once
)

func GetVersion() string {
	// getting the version of line-bot-sdk-go should be done only once. Computing it repeatedly is meaningless.
	sdkVersionOnce.Do(func() {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}
		for _, dep := range info.Deps {
			fmt.Print("dep.Path: ", dep.Path, " dep.Version: ", dep.Version, "\n")
			if strings.Contains(dep.Path, "github.com/kkdai/iloveptt") {
				sdkVersion = strings.TrimPrefix(dep.Version, "v")
				break
			}
		}
	})
	return sdkVersion
}
