package collector

import (
	"github.com/go-kit/kit/log"

	"github.com/prometheus/client_golang/prometheus"
	"time"
	"github.com/rspier/go-ecobee/ecobee"
	"fmt"
)

type EcobeeClient struct {
	APIKey       string   `toml:"apikey"`
	AuthCache       string   `toml:"auth_cache"`
	RefreshInterval       int64   `toml:"refresh_interval"`

}

type ecobeeCollector struct {
	fooMetric *prometheus.Desc
	barMetric *prometheus.Desc
	logger log.Logger
}

func init() {
	registerCollector("ecobee", defaultEnabled, newEcobeeCollector)
}

//You must create a constructor for you collector that
//initializes every descriptor and returns a pointer to the collector
func newEcobeeCollector(logger log.Logger) (Collector, error) {
	return &ecobeeCollector{
		fooMetric: prometheus.NewDesc("ecobee_metric",
			"Shows whether a foo has occurred in our cluster",
			nil, nil,
		),
		barMetric: prometheus.NewDesc("ecobee_bar_metric",
			"Shows whether a bar has occurred in our cluster",
			nil, nil,
		),
		logger: logger,
	}, nil
}

func (c *EcobeeClient) newEcobeeClient() (*ecobee.Client) {
	ecli := ecobee.NewClient(c.APIKey, c.AuthCache)

	return ecli
}

func (c *ecobeeCollector) Update(ch chan<- prometheus.Metric) error {
	ecli := newecobeeclien()
	if err != nil {
		return err
	}
	//Implement logic here to determine proper metric value to return to prometheus
	//for each descriptor or call other functions that do so.
	var metricValue float64
	if 1 == 1 {
		metricValue = 1
	}

	go func() {
		c := time.Tick(c.)
		for range c {
			refreshData(cli, *flagThermostatID)
		}
	}()


	//Write latest value for each metric in the prometheus metric channel.
	//Note that you can pass CounterValue, GaugeValue, or UntypedValue types here.
	ch <- prometheus.MustNewConstMetric(c.fooMetric, prometheus.CounterValue, metricValue)
	ch <- prometheus.MustNewConstMetric(c.barMetric, prometheus.CounterValue, metricValue)

	return nil
}
