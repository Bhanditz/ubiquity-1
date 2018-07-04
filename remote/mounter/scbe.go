/**
 * Copyright 2017 IBM Corp.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mounter

import (
	"fmt"
	"github.com/IBM/ubiquity/remote/mounter/block_device_mounter_utils"
	"github.com/IBM/ubiquity/resources"
	"github.com/IBM/ubiquity/utils"
	"github.com/IBM/ubiquity/utils/logs"
)

type scbeMounter struct {
	logger                  logs.Logger
	blockDeviceMounterUtils block_device_mounter_utils.BlockDeviceMounterUtils
	exec                    utils.Executor
	config                  resources.ScbeRemoteConfig
}

func NewScbeMounter(scbeRemoteConfig resources.ScbeRemoteConfig) resources.Mounter {
	blockDeviceMounterUtils := block_device_mounter_utils.NewBlockDeviceMounterUtils()
	return NewScbeMounterWithMounterUtil(scbeRemoteConfig, blockDeviceMounterUtils)
}

func NewScbeMounterWithMounterUtil(scbeRemoteConfig resources.ScbeRemoteConfig, blockDeviceMounterUtils block_device_mounter_utils.BlockDeviceMounterUtils) resources.Mounter {
	return &scbeMounter{
		logger:                  logs.GetLogger(),
		blockDeviceMounterUtils: blockDeviceMounterUtils,
		exec:   utils.NewExecutor(),
		config: scbeRemoteConfig,
	}
}

func (s *scbeMounter) Mount(mountRequest resources.MountRequest) (string, error) {
	defer s.logger.Trace(logs.DEBUG)()
	volumeWWN := mountRequest.VolumeConfig["Wwn"].(string)

	// Rescan OS
	if err := s.blockDeviceMounterUtils.RescanAll(!s.config.SkipRescanISCSI, volumeWWN, false); err != nil {
		return "", s.logger.ErrorRet(err, "RescanAll failed")
	}

	// Discover device
	devicePath, err := s.blockDeviceMounterUtils.Discover(volumeWWN, true)
	if err != nil {
		// Known issue: UB-1103 in https://www.ibm.com/support/knowledgecenter/SS6JWS_3.4.0/RN/sc_rn_knownissues.html
		// For DS8k, "rescan-scsi-bus.sh -r" cannot discover the LUN0, need to use "rescan-scsi-bus.sh -a" instead
		s.logger.Debug("volumeConfig: ", logs.Args{{"volumeConfig: ", mountRequest.VolumeConfig}})
		if isDS8kLun0(mountRequest) {
			if err := s.blockDeviceMounterUtils.RescanAll(!s.config.SkipRescanISCSI, volumeWWN, false); err != nil {
				return "", s.logger.ErrorRet(err, "RescanAll Targets failed", logs.Args{{"volumeWWN", volumeWWN}})
			}
			devicePath, err = s.blockDeviceMounterUtils.Discover(volumeWWN, true)
			if err != nil {
				return "", s.logger.ErrorRet(err, "Discover failed after rescan all target", logs.Args{{"volumeWWN", volumeWWN}})
			}
		} else {
			return "", s.logger.ErrorRet(err, "Discover failed", logs.Args{{"volumeWWN", volumeWWN}})
		}
	}

	// Create mount point if needed   // TODO consider to move it inside the util
	if _, err := s.exec.Stat(mountRequest.Mountpoint); err != nil {
		s.logger.Info("Create mountpoint directory " + mountRequest.Mountpoint)
		if err := s.exec.MkdirAll(mountRequest.Mountpoint, 0700); err != nil {
			return "", s.logger.ErrorRet(err, "MkdirAll failed", logs.Args{{"mountpoint", mountRequest.Mountpoint}})
		}
	}

	// Mount device and mkfs if needed
	var fstype string
	fstypeInterface, ok := mountRequest.VolumeConfig[resources.OptionNameForVolumeFsType]
	if !ok {
		// the backend should do this default, but this is just for safe
		fstype = resources.DefaultForScbeConfigParamDefaultFilesystem
	} else {
		fstype = fstypeInterface.(string)
	}

	if err := s.blockDeviceMounterUtils.MountDeviceFlow(devicePath, fstype, mountRequest.Mountpoint); err != nil {
		return "", s.logger.ErrorRet(err, "MountDeviceFlow failed", logs.Args{{"devicePath", devicePath}})
	}

	return mountRequest.Mountpoint, nil
}

func (s *scbeMounter) Unmount(unmountRequest resources.UnmountRequest) error {
	defer s.logger.Trace(logs.DEBUG)()

	volumeWWN := unmountRequest.VolumeConfig["Wwn"].(string)
	mountpoint := fmt.Sprintf(resources.PathToMountUbiquityBlockDevices, volumeWWN) // TODO instead of build the mountpoint it should come from unmountRequest.
	devicePath, err := s.blockDeviceMounterUtils.Discover(volumeWWN, true)
	if err != nil {
		// TODO idempotent, if volumeNotFoundError then skip UnmountDeviceFlow
		return s.logger.ErrorRet(err, "Discover failed", logs.Args{{"volumeWWN", volumeWWN}})
	}

	if err := s.blockDeviceMounterUtils.UnmountDeviceFlow(devicePath); err != nil {
		return s.logger.ErrorRet(err, "UnmountDeviceFlow failed", logs.Args{{"devicePath", devicePath}})
	}

	s.logger.Info("Delete mountpoint directory if exist", logs.Args{{"mountpoint", mountpoint}})
	// TODO move this part to the util
	if _, err := s.exec.Stat(mountpoint); err == nil {
		if err := s.exec.RemoveAll(mountpoint); err != nil { // TODO its enough to do Remove without All.
			return s.logger.ErrorRet(err, "RemoveAll failed", logs.Args{{"mountpoint", mountpoint}})
		}
	}

	return nil

}

func (s *scbeMounter) ActionAfterDetach(request resources.AfterDetachRequest) error {
	defer s.logger.Trace(logs.DEBUG)()
	volumeWWN := request.VolumeConfig["Wwn"].(string)

	// Rescan OS
	if err := s.blockDeviceMounterUtils.RescanAll(!s.config.SkipRescanISCSI, volumeWWN, true); err != nil {
		return s.logger.ErrorRet(err, "RescanAll failed")
	}
	return nil
}

func isDS8kLun0(mountRequest resources.MountRequest) bool {
	storageType, ok := mountRequest.VolumeConfig[resources.StorageType]
	if !ok {
		return false
	}
	lunNumber, ok := mountRequest.VolumeConfig[resources.ScbeKeyVolAttachLunNumToHost]
	if !ok {
		return false
	}

	if storageType == resources.DS8kStorageType && int(lunNumber.(float64)) == 0 {
		return true
	}
	return false
}