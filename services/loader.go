package services

import (
	"sync"
	"time"

	"github.com/shirou/gopsutil/cpu"
)

const (
	// reading cpu stat is a consuming operation, so we should establish some interval
	checkInterval = 60
	// during this time we calculate the load
	loadInterval = 1
)

type Loader interface {
	GetCpuUsage() (float64, error)
}

type LoaderStruct struct {
	cachedValue float64
	lastCheck   int64
	// as long as it is a consuming operation and in most cases we use cache, mutex looks like a good solution
	mu *sync.Mutex
}

func NewLoaderStruct() *LoaderStruct {
	return &LoaderStruct{
		cachedValue: defaultDifficulty,
		lastCheck:   time.Now().UnixNano(),
		mu:          &sync.Mutex{},
	}
}

// GetCpuUsage returns usage of CPU for calculating optimal load
func (l *LoaderStruct) GetCpuUsage() (float64, error) {
	if l.cachedValue != 0 && time.Now().UnixNano()-l.lastCheck > checkInterval {
		return l.cachedValue, nil
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	// we make extra check in order not to set lock before first check for optimisation.
	// so if several routines were locked and released, they don't have to recalculate usage
	if l.cachedValue != 0 && time.Now().UnixNano()-l.lastCheck > checkInterval {
		return l.cachedValue, nil
	}

	percentages, err := cpu.Percent(loadInterval, false)
	if err != nil {
		return 0, err
	}
	if len(percentages) == 0 {
		return 0, nil
	}

	l.lastCheck = time.Now().UnixNano()
	l.cachedValue = percentages[0]

	return l.cachedValue, nil
}
