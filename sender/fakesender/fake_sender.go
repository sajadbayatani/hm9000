// This file was generated by counterfeiter
package fakesender

import (
	"sync"

	"github.com/cloudfoundry/hm9000/models"
	"github.com/cloudfoundry/hm9000/sender"
	"code.cloudfoundry.org/clock"
)

type FakeSender struct {
	SendStub        func(clock.Clock, map[string]*models.App, []models.PendingStartMessage, []models.PendingStopMessage) error
	sendMutex       sync.RWMutex
	sendArgsForCall []struct {
		arg1 clock.Clock
		arg2 map[string]*models.App
		arg3 []models.PendingStartMessage
		arg4 []models.PendingStopMessage
	}
	sendReturns struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeSender) Send(arg1 clock.Clock, arg2 map[string]*models.App, arg3 []models.PendingStartMessage, arg4 []models.PendingStopMessage) error {
	var arg3Copy []models.PendingStartMessage
	if arg3 != nil {
		arg3Copy = make([]models.PendingStartMessage, len(arg3))
		copy(arg3Copy, arg3)
	}
	var arg4Copy []models.PendingStopMessage
	if arg4 != nil {
		arg4Copy = make([]models.PendingStopMessage, len(arg4))
		copy(arg4Copy, arg4)
	}
	fake.sendMutex.Lock()
	fake.sendArgsForCall = append(fake.sendArgsForCall, struct {
		arg1 clock.Clock
		arg2 map[string]*models.App
		arg3 []models.PendingStartMessage
		arg4 []models.PendingStopMessage
	}{arg1, arg2, arg3Copy, arg4Copy})
	fake.recordInvocation("Send", []interface{}{arg1, arg2, arg3Copy, arg4Copy})
	fake.sendMutex.Unlock()
	if fake.SendStub != nil {
		return fake.SendStub(arg1, arg2, arg3, arg4)
	} else {
		return fake.sendReturns.result1
	}
}

func (fake *FakeSender) SendCallCount() int {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return len(fake.sendArgsForCall)
}

func (fake *FakeSender) SendArgsForCall(i int) (clock.Clock, map[string]*models.App, []models.PendingStartMessage, []models.PendingStopMessage) {
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return fake.sendArgsForCall[i].arg1, fake.sendArgsForCall[i].arg2, fake.sendArgsForCall[i].arg3, fake.sendArgsForCall[i].arg4
}

func (fake *FakeSender) SendReturns(result1 error) {
	fake.SendStub = nil
	fake.sendReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeSender) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.sendMutex.RLock()
	defer fake.sendMutex.RUnlock()
	return fake.invocations
}

func (fake *FakeSender) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ sender.Sender = new(FakeSender)
