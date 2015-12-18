// This file was generated by counterfeiter
package fakes

import (
	"sync"

	"github.com/totherme/grufflo/backend"
	"github.com/totherme/grufflo/types"
)

type FakeParser struct {
	ParseStub        func(filePath string) (*types.GinkgoFile, error)
	parseMutex       sync.RWMutex
	parseArgsForCall []struct {
		filePath string
	}
	parseReturns struct {
		result1 *types.GinkgoFile
		result2 error
	}
}

func (fake *FakeParser) Parse(filePath string) (*types.GinkgoFile, error) {
	fake.parseMutex.Lock()
	fake.parseArgsForCall = append(fake.parseArgsForCall, struct {
		filePath string
	}{filePath})
	fake.parseMutex.Unlock()
	if fake.ParseStub != nil {
		return fake.ParseStub(filePath)
	} else {
		return fake.parseReturns.result1, fake.parseReturns.result2
	}
}

func (fake *FakeParser) ParseCallCount() int {
	fake.parseMutex.RLock()
	defer fake.parseMutex.RUnlock()
	return len(fake.parseArgsForCall)
}

func (fake *FakeParser) ParseArgsForCall(i int) string {
	fake.parseMutex.RLock()
	defer fake.parseMutex.RUnlock()
	return fake.parseArgsForCall[i].filePath
}

func (fake *FakeParser) ParseReturns(result1 *types.GinkgoFile, result2 error) {
	fake.ParseStub = nil
	fake.parseReturns = struct {
		result1 *types.GinkgoFile
		result2 error
	}{result1, result2}
}

var _ backend.Parser = new(FakeParser)