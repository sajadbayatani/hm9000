package md_test

import (
	"github.com/cloudfoundry/hm9000/models"
	"github.com/cloudfoundry/hm9000/testhelpers/app"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Stopping Duplicate Instances", func() {
	var a app.App
	var timestamp int

	Context("when there are multiple instances on the same index", func() {
		var instance0, instance1, duplicateInstance1 app.Instance
		var heartbeats []models.Heartbeat
		BeforeEach(func() {
			timestamp = 100
			a = app.NewApp()

			instance0 = a.GetInstance(0)
			instance1 = a.GetInstance(1)
			duplicateInstance1 = a.GetInstance(1)
			duplicateInstance1.InstanceGuid = app.Guid()

			heartbeats = []models.Heartbeat{models.Heartbeat{
				DeaGuid:            "abc",
				InstanceHeartbeats: []models.InstanceHeartbeat{instance0.Heartbeat(0), instance1.Heartbeat(0), duplicateInstance1.Heartbeat(0)},
			}}

			desired := a.DesiredState(0)
			desired.NumberOfInstances = 2
			stateServer.SetDesiredState([]models.DesiredAppState{desired})

			timestamp = sendHeartbeats(timestamp, heartbeats, 3, 10)
			cliRunner.Run("fetch_desired", timestamp)

			cliRunner.Run("analyze", timestamp)
		})

		It("should not immediately stop anything", func() {
			cliRunner.Run("send", timestamp)
			Ω(startStopListener.Stops).Should(BeEmpty())
		})

		Context("after a grace period", func() {
			BeforeEach(func() {
				timestamp += conf.GracePeriod
			})

			Context("if both instances are still running", func() {
				BeforeEach(func() {
					sendHeartbeats(timestamp, heartbeats, 1, 10)
					cliRunner.Run("analyze", timestamp)
					cliRunner.Run("send", timestamp)
				})

				It("should stop one of them", func() {
					Ω(startStopListener.Stops).Should(HaveLen(1))
					stop := startStopListener.Stops[0]
					Ω(stop.AppGuid).Should(Equal(a.AppGuid))
					Ω(stop.AppVersion).Should(Equal(a.AppVersion))
					Ω(stop.InstanceIndex).Should(Equal(1))
					Ω(stop.RunningIndices).Should(Equal(models.RunningIndices{"0": 1, "1": 2}))
					Ω([]string{instance1.InstanceGuid, duplicateInstance1.InstanceGuid}).Should(ContainElement(stop.InstanceGuid))
				})

				Context("after another grace period (assuming the stopped instance stops)", func() {
					BeforeEach(func() {
						timestamp += conf.GracePeriod
						instanceGuidToStop := startStopListener.Stops[0].InstanceGuid
						remainingInstance := instance1
						if remainingInstance.InstanceGuid == instanceGuidToStop {
							remainingInstance = duplicateInstance1
						}

						heartbeats = []models.Heartbeat{models.Heartbeat{
							DeaGuid:            "abc",
							InstanceHeartbeats: []models.InstanceHeartbeat{instance0.Heartbeat(0), remainingInstance.Heartbeat(0)},
						}}
						sendHeartbeats(timestamp, heartbeats, 1, 10)
						startStopListener.Reset()
						cliRunner.Run("analyze", timestamp)
						cliRunner.Run("send", timestamp)
					})

					It("should not stop the other instance", func() {
						Ω(startStopListener.Stops).Should(BeEmpty())
					})
				})
			})

			Context("if only one instance is still running", func() {
				BeforeEach(func() {
					heartbeats = []models.Heartbeat{models.Heartbeat{
						DeaGuid:            "abc",
						InstanceHeartbeats: []models.InstanceHeartbeat{instance0.Heartbeat(0), instance1.Heartbeat(0)},
					}}
					sendHeartbeats(timestamp, heartbeats, 1, 10)
					cliRunner.Run("analyze", timestamp)
					cliRunner.Run("send", timestamp)
				})

				It("should not stop any instances", func() {
					Ω(startStopListener.Stops).Should(BeEmpty())
				})
			})
		})
	})
})
