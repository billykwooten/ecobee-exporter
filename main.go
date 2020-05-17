// Copyright 2020 Billy Wooten
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

	"github.com/billykwooten/ecobee-exporter/collector"
	"github.com/billykwooten/go-ecobee/ecobee"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	app            = kingpin.New("ecobee-exporter", "Ecobee Exporter utilizing Ecobee API").Author("Billy Wooten")
	addr           = app.Flag("listen-address", "HTTP port to listen on").Envar("ECOBEE_LISTEN_ADDRESS").Default(":9098").String()
	applicationKey = app.Flag("appkey", "Application API Key").Envar("ECOBEE_APPKEY").Required().String()
	cacheFile      = app.Flag("cachefile", "Cache file so the exporter can store and sync authorization tokens").Envar("ECOBEE_CACHEFILE").Default("/db/auth.cache").String()
)

func main() {
	// Parse Kingpin Variables
	kingpin.MustParse(app.Parse(os.Args[1:]))

	// Setup Scopes for API Requests
	ecobee.Scopes = []string{"smartRead"}

	//Create a new instance of the ecobeeCollector and
	//register it with the prometheus client.
	ecobeeCollector := collector.EcobeeCollector(ecobee.NewClient(*applicationKey, *cacheFile), "ecobee")
	prometheus.MustRegister(ecobeeCollector)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port " + *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
