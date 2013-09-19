package store

import (
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/config"
	. "github.com/onsi/gomega"

	"github.com/cloudfoundry/hm9000/testhelpers/storerunner"

	"os"
	"os/signal"
	"testing"
)

var etcdRunner *storerunner.ETCDClusterRunner
var zookeeperRunner *storerunner.ZookeeperClusterRunner

func TestStore(t *testing.T) {
	registerSignalHandler()
	RegisterFailHandler(Fail)

	etcdPort := 5000 + (config.GinkgoConfig.ParallelNode-1)*10
	etcdRunner = storerunner.NewETCDClusterRunner(etcdPort, 5)

	zookeeperPort := 2181 + (config.GinkgoConfig.ParallelNode-1)*10
	zookeeperRunner = storerunner.NewZookeeperClusterRunner(zookeeperPort, 1)

	etcdRunner.Start()
	zookeeperRunner.Start()

	RunSpecs(t, "Store Suite")

	stopStores()
}

var _ = BeforeEach(func() {
	if etcdRunner != nil {
		etcdRunner.Reset()
	}

	if zookeeperRunner != nil {
		zookeeperRunner.Reset()
	}
})

func stopStores() {
	if etcdRunner != nil {
		etcdRunner.Stop()
	}

	if zookeeperRunner != nil {
		zookeeperRunner.Stop()
	}
}

func registerSignalHandler() {
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, os.Kill)

		select {
		case <-c:
			stopStores()
			os.Exit(0)
		}
	}()
}