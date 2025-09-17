package main

import (
	"fmt"
	"go-template/internal/app"
	"go-template/internal/config"
	"go-template/internal/log"
	"runtime/debug"

	"github.com/KimMachineGun/automemlimit/memlimit"
	"github.com/pbnjay/memory"

	"go.uber.org/automaxprocs/maxprocs"
)

func main() {
	log.Init()

	// Configure runtime.GOMAXPROCS to match the configured Linux CPU quota
	// The GOMAXPROCS variable limits the number of operating system threads that can execute user-level Go code simultaneously
	undo, err := maxprocs.Set(maxprocs.Logger(log.StdInfo))
	defer undo()

	if err != nil {
		log.Warnw("Failed to set GOMAXPROCS", "error", err)
	}

	// Configure GOMEMLIMIT
	// The GOMEMLIMIT variable sets a soft memory limit for the runtime.
	// This memory limit includes the Go heap and all other memory managed by the runtime,
	// and excludes external memory sources such as mappings of the binary itself, memory managed in other languages,
	// and memory held by the operating system on behalf of the Go program
	sysTotalMem := memory.TotalMemory()
	limit, err := memlimit.FromCgroup()

	if err == nil && limit < sysTotalMem {
		mem, _ := memlimit.SetGoMemLimit(0.9)
		log.Infow("Set memory limit by cgroup", "mem", memByteToStr(mem), "system_total_mem", memByteToStr(sysTotalMem))
	} else {
		mem := int64(float64(sysTotalMem) * 0.9)
		debug.SetMemoryLimit(mem)
		log.Infow("Set memory limit by system total memory", "mem", memByteToStr(mem), "system_total_mem", memByteToStr(sysTotalMem))
	}

	config.LoadEnv()

	// takes arguments to run migration actions
	handleArguments()

	// setup graceful reload
	signal := 2

	for signal >= 1 {
		signal = app.Run(config.GetEnv())
	}
}

func memByteToStr[T int64 | uint64](v T) string {
	return fmt.Sprintf("%d MB", uint64(v)/1048576)
}
