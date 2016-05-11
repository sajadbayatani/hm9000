package hm

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/cloudfoundry-incubator/consuladapter"
	"github.com/cloudfoundry-incubator/locket"
	"github.com/cloudfoundry/hm9000/analyzer"
	"github.com/cloudfoundry/hm9000/config"
	"github.com/cloudfoundry/hm9000/desiredstatefetcher"
	"github.com/cloudfoundry/hm9000/helpers/httpclient"
	"github.com/cloudfoundry/hm9000/helpers/metricsaccountant"
	"github.com/cloudfoundry/hm9000/models"
	"github.com/cloudfoundry/hm9000/sender"
	"github.com/cloudfoundry/hm9000/store"
	"github.com/cloudfoundry/yagnats"
	"github.com/pivotal-golang/clock"
	"github.com/pivotal-golang/lager"
)

func Analyze(l lager.Logger, sink lager.Sink, conf *config.Config, poll bool) {
	store := connectToStore(l, conf)
	messageBus := connectToMessageBus(l, conf)
	clock := buildClock(l)
	client := httpclient.NewHttpClient(conf.SkipSSLVerification, conf.FetcherNetworkTimeout())

	if poll {
		l.Info("Starting Analyzer...")

		f := &Component{
			component:       "analyzer",
			conf:            conf,
			pollingInterval: conf.FetcherPollingInterval(),
			timeout:         conf.FetcherTimeout(),
			logger:          l,
			action: func() error {
				return analyze(l, sink, clock, client, conf, store, messageBus)
			},
		}

		consulClient, _ := consuladapter.NewClientFromUrl(conf.ConsulCluster)
		lockRunner := locket.NewLock(l, consulClient, "hm9000.analyzer", make([]byte, 0), clock, locket.RetryInterval, locket.LockTTL)

		err := ifritizeComponent(f, lockRunner)

		if err != nil {
			l.Error("Analyzer exited", err)
			os.Exit(197)
		}

		l.Info("exited")
		os.Exit(0)
	} else {
		err := analyze(l, sink, clock, client, conf, store, messageBus)
		if err != nil {
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}

func analyze(l lager.Logger, sink lager.Sink, clock clock.Clock, client httpclient.HttpClient, conf *config.Config, store store.Store, messageBus yagnats.NATSConn) error {
	logger := lager.NewLogger("fetcher")
	logger.RegisterSink(sink)

	appQueue := make(chan map[string]models.DesiredAppState, 5)

	fetchDesiredErr := make(chan error)
	go func() {
		e := fetchDesiredState(logger, clock, client, conf, appQueue)
		fetchDesiredErr <- e
	}()

	analyzeStateErr := make(chan error)
	analyzeStateApps := make(chan map[string]*models.App)

	go func() {
		apps, e := analyzeState(l, clock, conf, store, appQueue)
		analyzeStateErr <- e
		analyzeStateApps <- apps
	}()

	var apps map[string]*models.App

ANALYZE_LOOP:
	for {
		select {
		case desiredErr := <-fetchDesiredErr:
			if desiredErr != nil {
				return desiredErr
			}
		case analyzeErr := <-analyzeStateErr:
			if analyzeErr != nil {
				return analyzeErr
			}
		case apps = <-analyzeStateApps:
			if apps != nil {
				break ANALYZE_LOOP
			}
		}
	}

	logger = lager.NewLogger("sender")
	logger.RegisterSink(sink)
	return send(logger, conf, messageBus, store, clock, apps)
}

func fetchDesiredState(l lager.Logger, clock clock.Clock, client httpclient.HttpClient,
	conf *config.Config, appQueue chan map[string]models.DesiredAppState) error {
	l.Info("Fetching Desired State")
	fetcher := desiredstatefetcher.New(
		conf,
		metricsaccountant.New(),
		client,
		clock,
		l,
	)

	resultChan := make(chan desiredstatefetcher.DesiredStateFetcherResult, 1)
	fetcher.Fetch(resultChan, appQueue)

	result := <-resultChan

	if result.Success {
		l.Info("Success", lager.Data{"Number of Desired Apps Fetched": strconv.Itoa(result.NumResults)})
		return nil
	}

	l.Error(result.Message, result.Error)
	return result.Error
}

func analyzeState(l lager.Logger, clk clock.Clock, conf *config.Config,
	store store.Store, appQueue chan map[string]models.DesiredAppState) (map[string]*models.App, error) {
	l.Info("Analyzing...")

	t := time.Now()
	a := analyzer.New(store, clk, l, conf)
	apps, err := a.Analyze(appQueue)
	analyzer.SendMetrics(apps, err)

	if err != nil {
		l.Error("Analyzer failed with error", err)
		return nil, err
	}

	l.Info("Analyzer completed succesfully", lager.Data{
		"Duration": fmt.Sprintf("%.4f", time.Since(t).Seconds()),
	})
	return apps, nil
}

func send(l lager.Logger, conf *config.Config, messageBus yagnats.NATSConn, store store.Store, clock clock.Clock, apps map[string]*models.App) error {
	l.Info("Sending...")

	sender := sender.New(store, metricsaccountant.New(), conf, messageBus, l, clock)
	err := sender.Send(clock, apps)

	if err != nil {
		l.Error("Sender failed with error", err)
		return err
	}

	l.Info("Sender completed succesfully")
	return nil
}
