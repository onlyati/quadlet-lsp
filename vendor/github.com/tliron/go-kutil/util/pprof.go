package util

import (
	"os"
	"runtime/pprof"
)

func InitializeCPUProfiling(path string) {
	if path == "" {
		return
	}

	cpuProfile, err := os.Create(path)
	FailOnError(err)
	OnExitError(cpuProfile.Close)
	err = pprof.StartCPUProfile(cpuProfile)
	FailOnError(err)
	OnExit(pprof.StopCPUProfile)
}
