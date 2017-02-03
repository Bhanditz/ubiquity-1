package spectrumscale

import (
	"fmt"
	"log"
	"os/exec"

	"github.ibm.com/almaden-containers/ubiquity/model"
	"github.ibm.com/almaden-containers/ubiquity/utils"
)

type spectrumNfsLocalClient struct {
	spectrumClient *spectrumLocalClient
	config         model.SpectrumScaleConfig
}

func NewSpectrumNfsLocalClient(logger *log.Logger, config model.SpectrumScaleConfig, dbClient utils.DatabaseClient, fileLock utils.FileLock) (model.StorageClient, error) {
	logger.Println("spectrumNfsLocalClient: init start")
	defer logger.Println("spectrumNfsLocalClient: init end")

	if config.ConfigPath == "" {
		return nil, fmt.Errorf("spectrumNfsLocalClient: init: missing required parameter 'spectrumConfigPath'")
	}

	if config.DefaultFilesystem == "" {
		return nil, fmt.Errorf("spectrumNfsLocalClient: init: missing required parameter 'spectrumDefaultFileSystem'")
	}

	if config.NfsServerAddr == "" {
		return nil, fmt.Errorf("spectrumNfsLocalClient: init: missing required parameter 'spectrumNfsServerAddr'")
	}

	spectrumClient, err := newSpectrumLocalClient(logger, config, dbClient, fileLock)
	if err != nil {
		return nil, err
	}
	return &spectrumNfsLocalClient{config: config, spectrumClient: spectrumClient}, nil
}

func (s *spectrumNfsLocalClient) Activate() error {
	s.spectrumClient.logger.Println("spectrumNfsLocalClient: Activate-start")
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: Activate-end")

	return s.spectrumClient.Activate()
}

func (s *spectrumNfsLocalClient) ListVolumes() ([]model.VolumeMetadata, error) {
	s.spectrumClient.logger.Println("spectrumNfsLocalClient: List-volumes-start")
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: List-volumes-end")

	return s.spectrumClient.ListVolumes()

}
func (s *spectrumNfsLocalClient) Attach(name string) (string, error) {
	s.spectrumClient.logger.Println("spectrumNfsLocalClient: Attach-start")
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: Attach-end")

	_, volumeConfig, err := s.GetVolume(name)
	if err != nil {
		return "", err
	}
	nfsShare, ok := volumeConfig["nfs_share"].(string)
	if !ok {
		err = fmt.Errorf("error getting NFS share info from volume config for volume %s", name)
		s.spectrumClient.logger.Println("spectrumNfsLocalClient: error: %v", err.Error())
		return "", err
	}
	return nfsShare, nil
}

func (s *spectrumNfsLocalClient) Detach(name string) error {
	s.spectrumClient.logger.Println("spectrumNfsLocalClient: Detach-start")
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: Detach-end")

	_, _, err := s.spectrumClient.GetVolume(name)

	if err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error in no-op detach for volume %s: %#v\n", name, err)
		return err
	}

	return nil
}

func (s *spectrumNfsLocalClient) GetPluginName() string {
	return "spectrum-scale-nfs"
}

func (s *spectrumNfsLocalClient) CreateVolume(name string, opts map[string]interface{}) error {

	s.spectrumClient.logger.Printf("spectrumNfsLocalClient: Create-start with name %s and opts %#v\n", name, opts)
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: Create-end")

	nfsClientCIDR, ok := opts["nfsClientCIDR"].(string)
	if !ok {
		errorMsg := "Cannot create volume (opts missing required parameter 'nfsClientCIDR')"
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: Create: Error: %s", errorMsg)
		return fmt.Errorf(errorMsg)
	}

	spectrumOpts := make(map[string]interface{})
	for k, v := range opts {
		if k != "nfsClientCIDR" {
			spectrumOpts[k] = v
		}
	}

	if err := s.spectrumClient.CreateVolume(name, spectrumOpts); err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error creating volume %#v\n", err)
		return err
	}

	if _, err := s.spectrumClient.Attach(name); err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error attaching volume %#v\n; deleting volume", err)
		s.spectrumClient.RemoveVolume(name, true)
		return err
	}

	if err := s.exportNfs(name, nfsClientCIDR); err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error exporting volume %#v\n; deleting volume", err)
		s.spectrumClient.RemoveVolume(name, true)
		return err
	}
	return nil
}

func (s *spectrumNfsLocalClient) RemoveVolume(name string, forceDelete bool) error {
	s.spectrumClient.logger.Println("spectrumNfsLocalClient: Remove-start")
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: Remove-end")

	if err := s.unexportNfs(name); err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: Could not unexport volume %s (error=%s)", name, err.Error())
	}
	if err := s.spectrumClient.Detach(name); err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: Could not detach volume %s (error=%s)", name, err.Error())
	}
	return s.spectrumClient.RemoveVolume(name, forceDelete)
}

func (s *spectrumNfsLocalClient) GetVolume(name string) (model.VolumeMetadata, map[string]interface{}, error) {
	s.spectrumClient.logger.Println("spectrumNfsLocalClient: GetV-start")
	defer s.spectrumClient.logger.Println("spectrumNfsLocalClient: Get-end")

	volumeMetadata, volumeConfig, err := s.spectrumClient.GetVolume(name)
	if err != nil {
		return volumeMetadata, volumeConfig, err
	}
	nfsShare := fmt.Sprintf("%s:%s", s.config.NfsServerAddr, volumeMetadata.Mountpoint)
	volumeConfig["nfs_share"] = nfsShare
	s.spectrumClient.logger.Printf("spectrumNfsLocalClient: GetVolume: Adding nfs_share %s to volume config for volume %s\n", nfsShare, name)
	return volumeMetadata, volumeConfig, nil
}

func (s *spectrumNfsLocalClient) exportNfs(name, clientCIDR string) error {
	s.spectrumClient.logger.Printf("spectrumNfsLocalClient: ExportNfs start with name=%#v and clientCIDR=%#v\n", name, clientCIDR)
	defer s.spectrumClient.logger.Printf("spectrumNfsLocalClient: ExportNfs end")

	existingVolume, exists, err := s.spectrumClient.dataModel.GetVolume(name)

	if err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: DbClient.GetVolume returned error %#v\n", err.Error())
		return err
	}
	if exists == false {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: volume %#s not found\n", err)
		return err
	}

	spectrumCommand := "/usr/lpp/mmfs/bin/mmnfs"
	volumeMountpoint, _, err := s.spectrumClient.getVolumeMountPoint(existingVolume)
	if err != nil {
		return err
	}

	args := []string{"export", "add", volumeMountpoint, "--client", fmt.Sprintf("%s(Access_Type=RW,Protocols=3:4,Squash=no_root_squash)", clientCIDR)}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error %#v ExportNfs output: %#v\n", err, output)
		return fmt.Errorf("Failed to export fileset via Nfs: %s", err.Error())
	}
	s.spectrumClient.logger.Printf("spectrumNfsLocalClient: ExportNfs output: %s\n", string(output))
	return nil
}

func (s *spectrumNfsLocalClient) unexportNfs(name string) error {
	s.spectrumClient.logger.Printf("spectrumNfsLocalClient: UnexportNfs start with name=%s\n", name)
	defer s.spectrumClient.logger.Printf("spectrumNfsLocalClient: ExportNfs end")

	existingVolume, exists, err := s.spectrumClient.dataModel.GetVolume(name)

	if err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error getting volume %#s \n", err)
		return err
	}
	if exists == false {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: volume %#s not found\n", err)
		return err
	}

	spectrumCommand := "/usr/lpp/mmfs/bin/mmnfs"
	volumeMountpoint, _, err := s.spectrumClient.getVolumeMountPoint(existingVolume)
	if err != nil {
		return err
	}

	args := []string{"export", "remove", volumeMountpoint, "--force"}
	cmd := exec.Command(spectrumCommand, args...)
	output, err := cmd.Output()
	if err != nil {
		s.spectrumClient.logger.Printf("spectrumNfsLocalClient: error %#v executing mmnfs command for output %#v \n", err, output)
		return fmt.Errorf("spectrumNfsLocalClient: Failed to unexport fileset via Nfs: %s", err.Error())
	}
	s.spectrumClient.logger.Printf("spectrumNfsLocalClient: UnexportNfs output: %s\n", string(output))

	return nil
}
