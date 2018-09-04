// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/IBM/ubiquity/remote/mounter/block_device_mounter_utils"
)

type FakeBlockDeviceMounterUtils struct {
	RescanAllStub        func(withISCSI bool, wwn string, rescanForCleanUp bool, ds8kLun0 bool) error
	rescanAllMutex       sync.RWMutex
	rescanAllArgsForCall []struct {
		withISCSI        bool
		wwn              string
		rescanForCleanUp bool
		ds8kLun0         bool
	}
	rescanAllReturns struct {
		result1 error
	}
	rescanAllReturnsOnCall map[int]struct {
		result1 error
	}
	MountDeviceFlowStub        func(devicePath string, fsType string, mountPoint string) error
	mountDeviceFlowMutex       sync.RWMutex
	mountDeviceFlowArgsForCall []struct {
		devicePath string
		fsType     string
		mountPoint string
	}
	mountDeviceFlowReturns struct {
		result1 error
	}
	mountDeviceFlowReturnsOnCall map[int]struct {
		result1 error
	}
	DiscoverStub        func(volumeWwn string, deepDiscovery bool) (string, error)
	discoverMutex       sync.RWMutex
	discoverArgsForCall []struct {
		volumeWwn     string
		deepDiscovery bool
	}
	discoverReturns struct {
		result1 string
		result2 error
	}
	discoverReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	UnmountDeviceFlowStub        func(devicePath string) error
	unmountDeviceFlowMutex       sync.RWMutex
	unmountDeviceFlowArgsForCall []struct {
		devicePath string
	}
	unmountDeviceFlowReturns struct {
		result1 error
	}
	unmountDeviceFlowReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeBlockDeviceMounterUtils) RescanAll(withISCSI bool, wwn string, rescanForCleanUp bool, ds8kLun0 bool) error {
	fake.rescanAllMutex.Lock()
	ret, specificReturn := fake.rescanAllReturnsOnCall[len(fake.rescanAllArgsForCall)]
	fake.rescanAllArgsForCall = append(fake.rescanAllArgsForCall, struct {
		withISCSI        bool
		wwn              string
		rescanForCleanUp bool
		ds8kLun0         bool
	}{withISCSI, wwn, rescanForCleanUp, ds8kLun0})
	fake.recordInvocation("RescanAll", []interface{}{withISCSI, wwn, rescanForCleanUp, ds8kLun0})
	fake.rescanAllMutex.Unlock()
	if fake.RescanAllStub != nil {
		return fake.RescanAllStub(withISCSI, wwn, rescanForCleanUp, ds8kLun0)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.rescanAllReturns.result1
}

func (fake *FakeBlockDeviceMounterUtils) RescanAllCallCount() int {
	fake.rescanAllMutex.RLock()
	defer fake.rescanAllMutex.RUnlock()
	return len(fake.rescanAllArgsForCall)
}

func (fake *FakeBlockDeviceMounterUtils) RescanAllArgsForCall(i int) (bool, string, bool, bool) {
	fake.rescanAllMutex.RLock()
	defer fake.rescanAllMutex.RUnlock()
	return fake.rescanAllArgsForCall[i].withISCSI, fake.rescanAllArgsForCall[i].wwn, fake.rescanAllArgsForCall[i].rescanForCleanUp, fake.rescanAllArgsForCall[i].ds8kLun0
}

func (fake *FakeBlockDeviceMounterUtils) RescanAllReturns(result1 error) {
	fake.RescanAllStub = nil
	fake.rescanAllReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockDeviceMounterUtils) RescanAllReturnsOnCall(i int, result1 error) {
	fake.RescanAllStub = nil
	if fake.rescanAllReturnsOnCall == nil {
		fake.rescanAllReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.rescanAllReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockDeviceMounterUtils) MountDeviceFlow(devicePath string, fsType string, mountPoint string) error {
	fake.mountDeviceFlowMutex.Lock()
	ret, specificReturn := fake.mountDeviceFlowReturnsOnCall[len(fake.mountDeviceFlowArgsForCall)]
	fake.mountDeviceFlowArgsForCall = append(fake.mountDeviceFlowArgsForCall, struct {
		devicePath string
		fsType     string
		mountPoint string
	}{devicePath, fsType, mountPoint})
	fake.recordInvocation("MountDeviceFlow", []interface{}{devicePath, fsType, mountPoint})
	fake.mountDeviceFlowMutex.Unlock()
	if fake.MountDeviceFlowStub != nil {
		return fake.MountDeviceFlowStub(devicePath, fsType, mountPoint)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.mountDeviceFlowReturns.result1
}

func (fake *FakeBlockDeviceMounterUtils) MountDeviceFlowCallCount() int {
	fake.mountDeviceFlowMutex.RLock()
	defer fake.mountDeviceFlowMutex.RUnlock()
	return len(fake.mountDeviceFlowArgsForCall)
}

func (fake *FakeBlockDeviceMounterUtils) MountDeviceFlowArgsForCall(i int) (string, string, string) {
	fake.mountDeviceFlowMutex.RLock()
	defer fake.mountDeviceFlowMutex.RUnlock()
	return fake.mountDeviceFlowArgsForCall[i].devicePath, fake.mountDeviceFlowArgsForCall[i].fsType, fake.mountDeviceFlowArgsForCall[i].mountPoint
}

func (fake *FakeBlockDeviceMounterUtils) MountDeviceFlowReturns(result1 error) {
	fake.MountDeviceFlowStub = nil
	fake.mountDeviceFlowReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockDeviceMounterUtils) MountDeviceFlowReturnsOnCall(i int, result1 error) {
	fake.MountDeviceFlowStub = nil
	if fake.mountDeviceFlowReturnsOnCall == nil {
		fake.mountDeviceFlowReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.mountDeviceFlowReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockDeviceMounterUtils) Discover(volumeWwn string, deepDiscovery bool) (string, error) {
	fake.discoverMutex.Lock()
	ret, specificReturn := fake.discoverReturnsOnCall[len(fake.discoverArgsForCall)]
	fake.discoverArgsForCall = append(fake.discoverArgsForCall, struct {
		volumeWwn     string
		deepDiscovery bool
	}{volumeWwn, deepDiscovery})
	fake.recordInvocation("Discover", []interface{}{volumeWwn, deepDiscovery})
	fake.discoverMutex.Unlock()
	if fake.DiscoverStub != nil {
		return fake.DiscoverStub(volumeWwn, deepDiscovery)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.discoverReturns.result1, fake.discoverReturns.result2
}

func (fake *FakeBlockDeviceMounterUtils) DiscoverCallCount() int {
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	return len(fake.discoverArgsForCall)
}

func (fake *FakeBlockDeviceMounterUtils) DiscoverArgsForCall(i int) (string, bool) {
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	return fake.discoverArgsForCall[i].volumeWwn, fake.discoverArgsForCall[i].deepDiscovery
}

func (fake *FakeBlockDeviceMounterUtils) DiscoverReturns(result1 string, result2 error) {
	fake.DiscoverStub = nil
	fake.discoverReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeBlockDeviceMounterUtils) DiscoverReturnsOnCall(i int, result1 string, result2 error) {
	fake.DiscoverStub = nil
	if fake.discoverReturnsOnCall == nil {
		fake.discoverReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.discoverReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *FakeBlockDeviceMounterUtils) UnmountDeviceFlow(devicePath string) error {
	fake.unmountDeviceFlowMutex.Lock()
	ret, specificReturn := fake.unmountDeviceFlowReturnsOnCall[len(fake.unmountDeviceFlowArgsForCall)]
	fake.unmountDeviceFlowArgsForCall = append(fake.unmountDeviceFlowArgsForCall, struct {
		devicePath string
	}{devicePath})
	fake.recordInvocation("UnmountDeviceFlow", []interface{}{devicePath})
	fake.unmountDeviceFlowMutex.Unlock()
	if fake.UnmountDeviceFlowStub != nil {
		return fake.UnmountDeviceFlowStub(devicePath)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.unmountDeviceFlowReturns.result1
}

func (fake *FakeBlockDeviceMounterUtils) UnmountDeviceFlowCallCount() int {
	fake.unmountDeviceFlowMutex.RLock()
	defer fake.unmountDeviceFlowMutex.RUnlock()
	return len(fake.unmountDeviceFlowArgsForCall)
}

func (fake *FakeBlockDeviceMounterUtils) UnmountDeviceFlowArgsForCall(i int) string {
	fake.unmountDeviceFlowMutex.RLock()
	defer fake.unmountDeviceFlowMutex.RUnlock()
	return fake.unmountDeviceFlowArgsForCall[i].devicePath
}

func (fake *FakeBlockDeviceMounterUtils) UnmountDeviceFlowReturns(result1 error) {
	fake.UnmountDeviceFlowStub = nil
	fake.unmountDeviceFlowReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockDeviceMounterUtils) UnmountDeviceFlowReturnsOnCall(i int, result1 error) {
	fake.UnmountDeviceFlowStub = nil
	if fake.unmountDeviceFlowReturnsOnCall == nil {
		fake.unmountDeviceFlowReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.unmountDeviceFlowReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeBlockDeviceMounterUtils) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.rescanAllMutex.RLock()
	defer fake.rescanAllMutex.RUnlock()
	fake.mountDeviceFlowMutex.RLock()
	defer fake.mountDeviceFlowMutex.RUnlock()
	fake.discoverMutex.RLock()
	defer fake.discoverMutex.RUnlock()
	fake.unmountDeviceFlowMutex.RLock()
	defer fake.unmountDeviceFlowMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeBlockDeviceMounterUtils) recordInvocation(key string, args []interface{}) {
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

var _ block_device_mounter_utils.BlockDeviceMounterUtils = new(FakeBlockDeviceMounterUtils)
