// This file was generated by counterfeiter
package fakemetricsaccountant

import (
	"sync"
	"time"

	"github.com/cloudfoundry/hm9000/helpers/metricsaccountant"
)

type FakeUsageTracker struct {
	StartTrackingUsageStub        func()
	startTrackingUsageMutex       sync.RWMutex
	startTrackingUsageArgsForCall []struct{}
	MeasureUsageStub              func() (usage float64, measurementDuration time.Duration)
	measureUsageMutex             sync.RWMutex
	measureUsageArgsForCall       []struct{}
	measureUsageReturns           struct {
		result1 float64
		result2 time.Duration
	}
}

func (fake *FakeUsageTracker) StartTrackingUsage() {
	fake.startTrackingUsageMutex.Lock()
	fake.startTrackingUsageArgsForCall = append(fake.startTrackingUsageArgsForCall, struct{}{})
	fake.startTrackingUsageMutex.Unlock()
	if fake.StartTrackingUsageStub != nil {
		fake.StartTrackingUsageStub()
	}
}

func (fake *FakeUsageTracker) StartTrackingUsageCallCount() int {
	fake.startTrackingUsageMutex.RLock()
	defer fake.startTrackingUsageMutex.RUnlock()
	return len(fake.startTrackingUsageArgsForCall)
}

func (fake *FakeUsageTracker) MeasureUsage() (usage float64, measurementDuration time.Duration) {
	fake.measureUsageMutex.Lock()
	fake.measureUsageArgsForCall = append(fake.measureUsageArgsForCall, struct{}{})
	fake.measureUsageMutex.Unlock()
	if fake.MeasureUsageStub != nil {
		return fake.MeasureUsageStub()
	} else {
		return fake.measureUsageReturns.result1, fake.measureUsageReturns.result2
	}
}

func (fake *FakeUsageTracker) MeasureUsageCallCount() int {
	fake.measureUsageMutex.RLock()
	defer fake.measureUsageMutex.RUnlock()
	return len(fake.measureUsageArgsForCall)
}

func (fake *FakeUsageTracker) MeasureUsageReturns(result1 float64, result2 time.Duration) {
	fake.MeasureUsageStub = nil
	fake.measureUsageReturns = struct {
		result1 float64
		result2 time.Duration
	}{result1, result2}
}

var _ metricsaccountant.UsageTracker = new(FakeUsageTracker)
