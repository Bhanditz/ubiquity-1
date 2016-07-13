package core_test

import (
	"fmt"
	"log"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gitlabhost.rtp.raleigh.ibm.com/ibm-storage/ibm-storage-broker/core"
	"gitlabhost.rtp.raleigh.ibm.com/ibm-storage/ibm-storage-broker/model"
	"gitlabhost.rtp.raleigh.ibm.com/ibm-storage/ibm-storage-broker/core/fakes"
)

var _ = Describe("ibm-storage-broker Broker", func() {
	var (
		controller      core.Controller
		localMountPoint string
		serviceGuid     string
		instanceMap     map[string]*model.ServiceInstance
		bindingMap      map[string]*model.ServiceBinding
		testLogger      log.Logger
		fakeBackend     *fakes.FakeStorageBackend
		configPath      string
	)
	BeforeEach(func() {
		serviceGuid = "some-service-guid"
		localMountPoint = "/tmp/share"
		configPath = "/tmp/ibm-storage-broker"
		instanceMap = make(map[string]*model.ServiceInstance)
		bindingMap = make(map[string]*model.ServiceBinding)
		fakeBackend = new(fakes.FakeStorageBackend)
		controller = core.NewController(fakeBackend, configPath, instanceMap, bindingMap)
	})
	Context(".Catalog", func() {
		It("should produce a valid catalog", func() {
			catalog, err := controller.GetCatalog(testLogger)
			Expect(err).ToNot(HaveOccurred())
			Expect(catalog).ToNot(BeNil())
			Expect(catalog.Services).ToNot(BeNil())
			Expect(len(catalog.Services)).To(Equal(1))
			Expect(catalog.Services[0].Name).To(Equal("ibm-storage-broker"))
			Expect(catalog.Services[0].Requires).ToNot(BeNil())
			Expect(len(catalog.Services[0].Requires)).To(Equal(1))
			Expect(catalog.Services[0].Requires[0]).To(Equal("volume_mount"))

			Expect(catalog.Services[0].Plans).ToNot(BeNil())
			Expect(len(catalog.Services[0].Plans)).To(Equal(1))
			Expect(catalog.Services[0].Plans[0].Name).To(Equal("free"))

			Expect(catalog.Services[0].Bindable).To(Equal(true))
		})
		Context(".CreateServiceInstance", func() {
			var (
				instance model.ServiceInstance
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}

			})
			It("should create a valid service instance", func() {
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
			})
			Context("should fail to create service instance", func() {
				It("when share creation errors", func() {
					properties := map[string]interface{}{"some-property": "some-value"}
					instance.Parameters = properties
					fakeBackend.CreateReturns(fmt.Errorf("Failed to create fileset"))
					_, err := controller.CreateServiceInstance(testLogger, "service-instance-guid", instance)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Failed to create fileset"))
				})
				It("should error when updating internal bookkeeping fails", func() {
					controller = core.NewController(fakeBackend, "/non-existent-path", instanceMap, bindingMap)
					_, err := controller.CreateServiceInstance(testLogger, "service-instance-guid", instance)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(fmt.Sprintf("open /non-existent-path/service_instances.json: no such file or directory")))
				})

			})
		})
		Context(".ServiceInstanceExists", func() {
			var (
				instance model.ServiceInstance
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}

			})
			It("should confirm existence of service instance", func() {
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				serviceExists := controller.ServiceInstanceExists(testLogger, serviceGuid)
				Expect(serviceExists).To(Equal(true))
			})
			It("should confirm non-existence of service instance", func() {
				serviceExists := controller.ServiceInstanceExists(testLogger, serviceGuid)
				Expect(serviceExists).To(Equal(false))
			})
		})
		Context(".ServiceInstancePropertiesMatch", func() {
			var (
				instance model.ServiceInstance
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}

			})
			It("should return true if properties match", func() {
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				anotherInstance := model.ServiceInstance{}
				properties := map[string]interface{}{"some-property": "some-value"}
				anotherInstance.Parameters = properties
				anotherInstance.PlanId = "some-planId"
				propertiesMatch := controller.ServiceInstancePropertiesMatch(testLogger, serviceGuid, anotherInstance)
				Expect(propertiesMatch).To(Equal(true))
			})
			It("should return false if properties do not match", func() {
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				anotherInstance := model.ServiceInstance{}
				properties := map[string]interface{}{"some-property": "some-value"}
				anotherInstance.Parameters = properties
				anotherInstance.PlanId = "some-other-planId"
				propertiesMatch := controller.ServiceInstancePropertiesMatch(testLogger, serviceGuid, anotherInstance)
				Expect(propertiesMatch).ToNot(Equal(true))
			})
		})

		Context(".ServiceInstanceDelete", func() {
			var (
				instance model.ServiceInstance
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}
			})
			It("should delete service instance", func() {
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				err := controller.DeleteServiceInstance(testLogger, serviceGuid)
				Expect(err).ToNot(HaveOccurred())

				serviceExists := controller.ServiceInstanceExists(testLogger, serviceGuid)
				Expect(serviceExists).To(Equal(false))
			})
			It("should error when trying to delete non-existence service instance", func() {
				fakeBackend.RemoveReturns(fmt.Errorf("Failed to delete fileset, fileset not found"))
				err := controller.DeleteServiceInstance(testLogger, serviceGuid)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal("Failed to delete fileset, fileset not found"))
			})
			It("should error when updating internal bookkeeping fails", func() {
				controller = core.NewController(fakeBackend, "/non-existent-path", instanceMap, bindingMap)
				err := controller.DeleteServiceInstance(testLogger, serviceGuid)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(Equal(fmt.Sprintf("open /non-existent-path/service_instances.json: no such file or directory")))
			})

		})
		Context(".BindServiceInstance", func() {
			var (
				instance    model.ServiceInstance
				bindingInfo model.ServiceBinding
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}
				bindingInfo = model.ServiceBinding{}
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
			})
			It("should be able bind service instance", func() {
				config := make(map[string]interface{})
				fakeBackend.GetReturns(&model.VolumeMetadata{Mountpoint: "/gpfs/fileset1"}, &config, nil)
				bindingResponse, err := controller.BindServiceInstance(testLogger, serviceGuid, "some-binding-id", bindingInfo)
				Expect(err).ToNot(HaveOccurred())
				Expect(bindingResponse.VolumeMounts).ToNot(BeNil())
				Expect(len(bindingResponse.VolumeMounts)).To(Equal(1))
			})
			Context("should fail", func() {
				It("when unable to find the backing share", func() {
					config := make(map[string]interface{})
					fakeBackend.GetReturns(&model.VolumeMetadata{}, &config, fmt.Errorf("Cannot find fileset, internal error"))
					_, err := controller.BindServiceInstance(testLogger, serviceGuid, "some-binding-id", bindingInfo)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal("Cannot find fileset, internal error"))
				})
				It("when updating internal bookkeeping fails", func() {
					config := make(map[string]interface{})
					fakeBackend.GetReturns(&model.VolumeMetadata{Mountpoint: "/gpfs/fileset1"}, &config, nil)
					controller = core.NewController(fakeBackend, "/non-existent-path", instanceMap, bindingMap)
					_, err := controller.BindServiceInstance(testLogger, serviceGuid, "some-binding-id", bindingInfo)
					Expect(err).To(HaveOccurred())
					Expect(err.Error()).To(Equal(fmt.Sprintf("open /non-existent-path/service_bindings.json: no such file or directory")))
				})
			})
		})
		Context(".ServiceBindingExists", func() {
			var (
				instance  model.ServiceInstance
				bindingId string
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}
				bindingId = "some-binding-id"
			})
			It("should confirm existence of service instance", func() {
				config := make(map[string]interface{})
				fakeBackend.GetReturns(&model.VolumeMetadata{Mountpoint: "/gpfs/fileset1"}, &config, nil)
				binding := model.ServiceBinding{}
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				successfulServiceBindingCreate(testLogger, fakeBackend, binding, controller, serviceGuid, bindingId)
				bindingExists := controller.ServiceBindingExists(testLogger, serviceGuid, bindingId)
				Expect(bindingExists).To(Equal(true))
			})
			It("should confirm non-existence of service binding", func() {
				bindingExists := controller.ServiceBindingExists(testLogger, serviceGuid, bindingId)
				Expect(bindingExists).To(Equal(false))
			})
		})
		Context(".ServiceBindingPropertiesMatch", func() {
			var (
				instance  model.ServiceInstance
				bindingId string
			)
			BeforeEach(func() {
				instance = model.ServiceInstance{}
				instance.PlanId = "some-planId"
				instance.Parameters = map[string]interface{}{"some-property": "some-value"}
				bindingId = "some-binding-id"

			})
			It("should return true if properties match", func() {
				binding := model.ServiceBinding{}
				config := make(map[string]interface{})
				fakeBackend.GetReturns(&model.VolumeMetadata{Mountpoint: "/gpfs/fileset1"}, &config, nil)
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				successfulServiceBindingCreate(testLogger, fakeBackend, binding, controller, serviceGuid, bindingId)
				anotherBinding := model.ServiceBinding{}
				propertiesMatch := controller.ServiceBindingPropertiesMatch(testLogger, serviceGuid, bindingId, anotherBinding)
				Expect(propertiesMatch).To(Equal(true))
			})
			It("should return false if properties do not match", func() {
				binding := model.ServiceBinding{}
				config := make(map[string]interface{})
				fakeBackend.GetReturns(&model.VolumeMetadata{Mountpoint: "/gpfs/fileset1"}, &config, nil)
				successfulServiceInstanceCreate(testLogger, fakeBackend, instance, controller, serviceGuid)
				successfulServiceBindingCreate(testLogger, fakeBackend, binding, controller, serviceGuid, bindingId)
				anotherBinding := model.ServiceBinding{}
				anotherBinding.AppId = "some-other-appId"
				propertiesMatch := controller.ServiceBindingPropertiesMatch(testLogger, serviceGuid, bindingId, anotherBinding)
				Expect(propertiesMatch).ToNot(Equal(true))
			})
		})
	})
})

func successfulServiceInstanceCreate(testLogger log.Logger, fakeBackend *fakes.FakeStorageBackend, instance model.ServiceInstance, controller core.Controller, serviceGuid string) {
	fakeBackend.CreateReturns(nil)
	createResponse, err := controller.CreateServiceInstance(testLogger, serviceGuid, instance)
	Expect(err).ToNot(HaveOccurred())
	Expect(createResponse.DashboardUrl).ToNot(Equal(""))
}

func successfulServiceBindingCreate(testLogger log.Logger, fakeBackend *fakes.FakeStorageBackend, binding model.ServiceBinding, controller core.Controller, serviceGuid string, bindingId string) {
	bindResponse, err := controller.BindServiceInstance(testLogger, serviceGuid, bindingId, binding)
	Expect(err).ToNot(HaveOccurred())
	Expect(bindResponse.VolumeMounts).ToNot(BeNil())
	Expect(len(bindResponse.VolumeMounts)).To(Equal(1))
}
