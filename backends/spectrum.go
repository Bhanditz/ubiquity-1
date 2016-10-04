package backends

import (
	"log"
	"fmt"
	"strings"

	common "github.ibm.com/almaden-containers/spectrum-common.git/core"
	"github.ibm.com/almaden-containers/ibm-storage-broker.git/model"
)

type SpectrumBackend struct {
	logger *log.Logger
	mountpoint string
	spectrumClients map[string]common.SpectrumClient  // cached SpectrumClient instance (one per service plan)
}

func NewSpectrumBackend(logger *log.Logger, mountpoint string) *SpectrumBackend {
	return &SpectrumBackend{logger: logger, mountpoint: mountpoint, spectrumClients: make(map[string]common.SpectrumClient)}
}

func (s *SpectrumBackend) GetServices() []model.Service {
	plan1 := model.ServicePlan{
		Name:        "gold",
		Id:          "spectrum-scale-gold",
		Description: "Gold Tier Performance and Resiliency",
		Metadata:    nil,
		Free:        true,
	}

	plan2 := model.ServicePlan{
		Name:        "bronze",
		Id:          "spectrum-scale-bronze",
		Description: "Bronze Tier Performance and Resiliency",
		Metadata:    nil,
		Free:        true,
	}

	service := model.Service{
		Name:            "spectrum-scale",
		Id:              "spectrum-service-guid",
		Description:     "Provides the Spectrum FS volume service, including volume creation and volume mounts",
		Bindable:        true,
		PlanUpdateable:  false,
		Tags:            []string{"gpfs"},
		Requires:        []string{"volume_mount"},
		Metadata:        nil,
		Plans:           []model.ServicePlan{plan1, plan2},
		DashboardClient: nil,
	}

	return []model.Service{service}
}

func (s *SpectrumBackend) CreateVolume(serviceInstance model.ServiceInstance, name string, opts map[string]interface{}) error {
	client := s.getSpectrumClient(serviceInstance)
	return client.Create(name, opts)
}

func (s *SpectrumBackend) RemoveVolume(serviceInstance model.ServiceInstance, name string) error {
	client := s.getSpectrumClient(serviceInstance)
	return client.Remove(name)
}

func (s *SpectrumBackend) ListVolumes(serviceInstance model.ServiceInstance) ([]model.VolumeMetadata, error){
	client := s.getSpectrumClient(serviceInstance)
	spectrumVolumeMetaData, err := client.List()

	volumeMetaData := make([]model.VolumeMetadata, len(spectrumVolumeMetaData))
	for i, e := range spectrumVolumeMetaData {
		volumeMetaData[i] = model.VolumeMetadata{
			Name: e.Name,
			Mountpoint: e.Mountpoint,
		}
	}

	return volumeMetaData, err
}

func (s *SpectrumBackend) GetVolume(serviceInstance model.ServiceInstance, name string) (volumeMetadata *model.VolumeMetadata, clientDriverName *string, config *map[string]interface{}, err error) {
	client := s.getSpectrumClient(serviceInstance)
	spectrumVolumeMetaData, spectrumConfig, err := client.Get(name)
	if err != nil {
		return nil, nil, nil, err
	}

	volumeMetadata = &model.VolumeMetadata{
		Name: spectrumVolumeMetaData.Name,
		Mountpoint: spectrumVolumeMetaData.Mountpoint,
	}

	configMap := make(map[string]interface{})
	configMap["fileset"] = spectrumConfig.FilesetId
	configMap["filesystem"] = spectrumConfig.Filesystem
	clientDriver := fmt.Sprintf("spectrum-scale-%s", spectrumConfig.Filesystem)

	return volumeMetadata, &clientDriver, &configMap, nil
}

func (s *SpectrumBackend) getSpectrumClient(serviceInstance model.ServiceInstance) common.SpectrumClient {
	// FIXME: clean up usage of planId for getting plan name
	planIdSplit := strings.Split(serviceInstance.PlanId, "-")
	planName := planIdSplit[len(planIdSplit)-1]

	// cache SpectrumClients per plan
	spectrumClient, ok := s.spectrumClients[planName]
	if !ok {
		mountPoint := fmt.Sprintf("%s/%s", s.mountpoint, planName)
		dbclient := common.NewDatabaseClient(s.logger, planName, mountPoint)
		dbclient.Init()
		spectrumClient := common.NewSpectrumClient(s.logger, planName, mountPoint, dbclient)
		spectrumClient.Activate()
		s.spectrumClients[planName] = spectrumClient
	}
	return spectrumClient
}
