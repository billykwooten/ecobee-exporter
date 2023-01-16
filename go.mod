module github.com/ejsuncy/ecobee-exporter

go 1.16

require (
	github.com/billykwooten/ecobee-exporter v0.0.2
	github.com/billykwooten/go-ecobee v0.0.1
	github.com/prometheus/client_golang v1.10.0
	github.com/sirupsen/logrus v1.8.1
	gopkg.in/alecthomas/kingpin.v2 v2.2.6
)

replace github.com/billykwooten/ecobee-exporter => github.com/ejsuncy/ecobee-exporter main-fork
