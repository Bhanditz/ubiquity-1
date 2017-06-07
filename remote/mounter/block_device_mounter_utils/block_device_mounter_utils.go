package block_device_mounter_utils

import (
	"fmt"
	"github.com/IBM/ubiquity/remote/mounter/block_device_utils"
)

// MountDeviceFlow create filesystem on the device (if needed) and then mount it on a given mountpoint
func (s *blockDeviceMounterUtils) MountDeviceFlow(devicePath string, fsType string, mountPoint string) error {
	s.logger.Printf(
		fmt.Sprintf("MountDevice: start to mount device [%s] (required with fstype [%s]) on mountpoint [%s]",
			devicePath, fsType, mountPoint))
	fsAlreadyExist, err := s.BlockDeviceUtilsInst.CheckFs(devicePath)
	if err != nil {
		return err
	}
	if !fsAlreadyExist {
		if err := s.BlockDeviceUtilsInst.MakeFs(devicePath, fsType); err != nil {
			return err
		}
	}
	if err := s.BlockDeviceUtilsInst.MountFs(devicePath, mountPoint); err != nil {
		return err
	}
	fmt.Sprintf("MountDevice: Successfully mount device ", devicePath)
	return nil
}

// RescanAll triggers the following OS rescanning :
// 1. iSCSI rescan (if protocol given is iscsi)
// 2. SCSI rescan
// 3. multipathing rescan
// return error if one of the steps fail
func (s *blockDeviceMounterUtils) RescanAll(withISCSI bool) error {
	s.logger.Printf("Start rescan OS i/SCSI devices and multipathing:")
	if withISCSI {
		if err := s.BlockDeviceUtilsInst.Rescan(block_device_utils.ISCSI); err != nil {
			return err
		}
	}
	if err := s.BlockDeviceUtilsInst.Rescan(block_device_utils.SCSI); err != nil {
		return err
	}

	if err := s.BlockDeviceUtilsInst.ReloadMultipath(); err != nil {
		return err
	}
	s.logger.Printf("Finished rescanning OS SCSI devices and multipathing:")
	return nil
}

func (s *blockDeviceMounterUtils) Discover(volumeWwn string) (string, error) {
	return s.BlockDeviceUtilsInst.Discover(volumeWwn)
}
