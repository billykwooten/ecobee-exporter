// Package prometheus provides Prometheus support for ecobee metrics.
package collector

import (
	"fmt"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/billykwooten/go-ecobee/ecobee"
	"github.com/prometheus/client_golang/prometheus"
)

type descs string

func (d descs) new(fqName, help string, variableLabels []string) *prometheus.Desc {
	return prometheus.NewDesc(fmt.Sprintf("%s_%s", d, fqName), help, variableLabels, nil)
}

// eCollector implements prometheus.eCollector to gather ecobee metrics on-demand.
type eCollector struct {
	client *ecobee.Client

	// per-query descriptors
	fetchTime *prometheus.Desc

	// runtime descriptors
	actualTemperature, targetTemperatureMin, targetTemperatureMax, lastIntervalRuntime, lastIntervalEnergizedStage *prometheus.Desc

	// sensor descriptors
	temperature, humidity, occupancy, inUse, currentHvacMode *prometheus.Desc

	// weather descriptors
	forecastTemperature, forecastCondition, forecastPressure, forecastRelativeHumidity, forecastDewpoint, forecastVisibility,
	forecastWindSpeed, forecastWindGust, forecastWindDirection, forecastWindBearing, forecastProbPrecip, forecastTempHigh,
	forecastTempLow, forecastSky *prometheus.Desc
}

// NewEcobeeCollector returns a new eCollector with the given prefix assigned to all
// metrics. Note that Prometheus metrics must be unique! Don't try to create
// two Collectors with the same metric prefix.
func NewEcobeeCollector(c *ecobee.Client, metricPrefix string) *eCollector {
	d := descs(metricPrefix)

	// fields common across multiple metrics
	runtime := []string{"thermostat_id", "thermostat_name"}
	sensor := append(runtime, "sensor_id", "sensor_name", "sensor_type")

	return &eCollector{
		client: c,

		// collector metrics
		fetchTime: d.new(
			"fetch_time",
			"elapsed time fetching data via Ecobee API",
			nil,
		),

		// thermostat (aka runtime) metrics
		actualTemperature: d.new(
			"actual_temperature",
			"thermostat-averaged current temperature",
			runtime,
		),
		targetTemperatureMax: d.new(
			"target_temperature_max",
			"maximum temperature for thermostat to maintain",
			runtime,
		),
		targetTemperatureMin: d.new(
			"target_temperature_min",
			"minimum temperature for thermostat to maintain",
			runtime,
		),
		lastIntervalRuntime: d.new(
			"last_interval_runtime",
			"last (latest reported) interval runtime of given stage mechanism (seconds, 0-300)",
			[]string{"thermostat_id", "thermostat_name", "stage_name", "stage_mechanism"},
		),
		lastIntervalEnergizedStage: d.new(
			"last_interval_energized_state",
			"last (latest reported) interval energized stage (one of: heatStage10n, heatStage20n, heatStage30n, heatOff, compressorCoolStage10n, compressorCoolStage20n, compressorCoolOff, compressorHeatStage10n, compressorHeatStage20n, compressorHeatOff, economyCycle",
			[]string{"thermostat_id", "thermostat_name", "stage"},
		),

		// sensor metrics
		temperature: d.new(
			"temperature",
			"temperature reported by a sensor in degrees",
			sensor,
		),
		humidity: d.new(
			"humidity",
			"humidity reported by a sensor in percent",
			sensor,
		),
		occupancy: d.new(
			"occupancy",
			"occupancy reported by a sensor (0 or 1)",
			sensor,
		),
		inUse: d.new(
			"in_use",
			"is sensor being used in thermostat calculations (0 or 1)",
			sensor,
		),
		currentHvacMode: d.new(
			"currenthvacmode",
			"current hvac mode of thermostat",
			[]string{"thermostat_id", "thermostat_name", "current_hvac_mode"},
		),

		// weather metrics
		forecastTemperature: d.new(
			"forecast_temperature",
			"weather forecast temperature for the thermostat",
			runtime,
		),
		forecastCondition: d.new(
			"forecast_condition",
			"weather forecast condition for the thermostat",
			[]string{"thermostat_id", "thermostat_name", "name"},
		),
		forecastPressure: d.new(
			"forecast_pressure",
			"weather forecast barometric pressure for the thermostat",
			runtime,
		),
		forecastRelativeHumidity: d.new(
			"forecast_relative_humidity",
			"weather forecast relative humidity for the thermostat (as a percent)",
			runtime,
		),
		forecastDewpoint: d.new(
			"forecast_dewpoint",
			"weather forecast dewpoint temperature for the thermostat",
			runtime,
		),
		forecastVisibility: d.new(
			"forecast_visibility",
			"weather forecast visibility for the thermostat (in meters, 0 - 70,000)",
			runtime,
		),
		forecastWindSpeed: d.new(
			"forecast_wind_speed",
			"weather forecast wind speed for the thermostat (in mph*1000)",
			runtime,
		),
		forecastWindGust: d.new(
			"forecast_wind_gust",
			"weather forecast wind gust for the thermostat (in mph*1000)",
			runtime,
		),
		forecastWindDirection: d.new(
			"forecast_wind_direction",
			"weather forecast wind direction for the thermostat",
			[]string{"thermostat_id", "thermostat_name", "direction"},
		),
		forecastWindBearing: d.new(
			"forecast_wind_bearing",
			"weather forecast wind bearing for the thermostat",
			runtime,
		),
		forecastProbPrecip: d.new(
			"forecast_probability_of_precipitation",
			"weather forecast probability of precipitation for the thermostat",
			runtime,
		),
		forecastTempHigh: d.new(
			"forecast_temp_high",
			"weather forecast high temperature for the day for the thermostat",
			runtime,
		),
		forecastTempLow: d.new(
			"forecast_temp_low",
			"weather forecast low temperature for the day for the thermostat",
			runtime,
		),
		forecastSky: d.new(
			"forecast_sky",
			"weather forecast sky condition for the thermostat",
			[]string{"thermostat_id", "thermostat_name", "condition"},
		),
	}
}

// Describe dumps all metric descriptors into ch.
func (c *eCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.fetchTime
	ch <- c.actualTemperature
	ch <- c.targetTemperatureMax
	ch <- c.targetTemperatureMin
	ch <- c.temperature
	ch <- c.humidity
	ch <- c.occupancy
	ch <- c.inUse
	ch <- c.currentHvacMode
	ch <- c.lastIntervalRuntime
}

// Collect retrieves thermostat data via the ecobee API.
func (c *eCollector) Collect(ch chan<- prometheus.Metric) {
	start := time.Now()
	tt, err := c.client.GetThermostats(ecobee.Selection{
		SelectionType:          "registered",
		IncludeSensors:         true,
		IncludeRuntime:         true,
		IncludeWeather:         true,
		IncludeSettings:        true,
		IncludeExtendedRuntime: true,
	})
	elapsed := time.Now().Sub(start)
	ch <- prometheus.MustNewConstMetric(c.fetchTime, prometheus.GaugeValue, elapsed.Seconds())
	if err != nil {
		log.Error(err)
		return
	}
	for _, t := range tt {
		tFields := []string{t.Identifier, t.Name}
		if t.Runtime.Connected {
			ch <- prometheus.MustNewConstMetric(
				c.actualTemperature, prometheus.GaugeValue, float64(t.Runtime.ActualTemperature)/10, tFields...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.targetTemperatureMax, prometheus.GaugeValue, float64(t.Runtime.DesiredCool)/10, tFields...,
			)
			ch <- prometheus.MustNewConstMetric(
				c.targetTemperatureMin, prometheus.GaugeValue, float64(t.Runtime.DesiredHeat)/10, tFields...,
			)

			// API returns last 3 intervals; we're only interested in the last (latest) one
			lastIndex := 2

			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.HeatPump1[lastIndex]), t.Identifier, t.Name, "heatStage1On", "heatPump1",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.HeatPump2[lastIndex]), t.Identifier, t.Name, "heatStage2On", "heatPump2",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.AuxHeat1[lastIndex]), t.Identifier, t.Name, "heatStage1On", "auxHeat1",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.AuxHeat2[lastIndex]), t.Identifier, t.Name, "heatStage2On", "auxHeat2",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Cool1[lastIndex]), t.Identifier, t.Name, "compressorCoolStage1On", "cool1",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Cool2[lastIndex]), t.Identifier, t.Name, "compressorCoolStage2On", "cool2",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Fan[lastIndex]), t.Identifier, t.Name, "fan", "fan",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Humidifier[lastIndex]), t.Identifier, t.Name, "humidifier", "humidifier",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Dehumidifier[lastIndex]), t.Identifier, t.Name, "dehumidifier", "dehumidifier",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Economizer[lastIndex]), t.Identifier, t.Name, "economizer", "economizer",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalRuntime, prometheus.GaugeValue, float64(t.ExtendedRuntime.Ventilator[lastIndex]), t.Identifier, t.Name, "ventilator", "ventilator",
			)
			ch <- prometheus.MustNewConstMetric(
				c.lastIntervalEnergizedStage, prometheus.GaugeValue, 0, t.Identifier, t.Name, t.ExtendedRuntime.HvacMode[lastIndex],
			)

			ch <- prometheus.MustNewConstMetric(
				c.currentHvacMode, prometheus.GaugeValue, 0, t.Identifier, t.Name, t.Settings.HvacMode,
			)

			// Weather
			// The first forecast is the most accurate per API docs
			forecast := t.Weather.Forecasts[0]

			ch <- prometheus.MustNewConstMetric(
				c.forecastTemperature, prometheus.GaugeValue, float64(forecast.Temperature)/10, tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastCondition, prometheus.GaugeValue, 1, t.Identifier, t.Name, forecast.Condition,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastPressure, prometheus.GaugeValue, float64(forecast.Pressure), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastRelativeHumidity, prometheus.GaugeValue, float64(forecast.RelativeHumidity), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastDewpoint, prometheus.GaugeValue, float64(forecast.Dewpoint)/10, tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastVisibility, prometheus.GaugeValue, float64(forecast.Visibility), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastWindSpeed, prometheus.GaugeValue, float64(forecast.WindSpeed), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastWindGust, prometheus.GaugeValue, float64(forecast.WindGust), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastWindDirection, prometheus.GaugeValue, 1, t.Identifier, t.Name, forecast.WindDirection,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastWindBearing, prometheus.GaugeValue, float64(forecast.WindBearing), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastProbPrecip, prometheus.GaugeValue, float64(forecast.Pop), tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastTempHigh, prometheus.GaugeValue, float64(forecast.TempHigh)/10, tFields...,
			)

			ch <- prometheus.MustNewConstMetric(
				c.forecastTempLow, prometheus.GaugeValue, float64(forecast.TempLow)/10, tFields...,
			)

			skyConditionStr := getSkyConditions()[int(forecast.Sky)]
			ch <- prometheus.MustNewConstMetric(
				c.forecastSky, prometheus.GaugeValue, 1, t.Identifier, t.Name, skyConditionStr,
			)

		}
		for _, s := range t.RemoteSensors {
			sFields := append(tFields, s.ID, s.Name, s.Type)
			inUse := float64(0)
			if s.InUse {
				inUse = 1
			}
			ch <- prometheus.MustNewConstMetric(
				c.inUse, prometheus.GaugeValue, inUse, sFields...,
			)
			for _, sc := range s.Capability {
				switch sc.Type {
				case "temperature":
					if v, err := strconv.ParseFloat(sc.Value, 64); err == nil {
						ch <- prometheus.MustNewConstMetric(
							c.temperature, prometheus.GaugeValue, v/10, sFields...,
						)
					} else {
						log.Error(err)
					}
				case "humidity":
					if v, err := strconv.ParseFloat(sc.Value, 64); err == nil {
						ch <- prometheus.MustNewConstMetric(
							c.humidity, prometheus.GaugeValue, v, sFields...,
						)
					} else {
						log.Error(err)
					}
				case "occupancy":
					switch sc.Value {
					case "true":
						ch <- prometheus.MustNewConstMetric(
							c.occupancy, prometheus.GaugeValue, 1, sFields...,
						)
					case "false":
						ch <- prometheus.MustNewConstMetric(
							c.occupancy, prometheus.GaugeValue, 0, sFields...,
						)
					default:
						log.Errorf("unknown sensor occupancy value %q", sc.Value)
					}
				default:
					log.Infof("ignoring sensor capability %q", sc.Type)
				}
			}
		}
	}
}

func getSkyConditions() []string {
	return []string{
		"UNDEFINED",
		"SUNNY",
		"CLEAR",
		"MOSTLY_SUNNY",
		"MOSTLY_CLEAR",
		"HAZY_SUNSHINE",
		"HAZE",
		"PASSING_CLOUDS",
		"MORE_SUN_THAN_CLOUDS",
		"SCATTERED_CLOUDS",
		"PARTLY_CLOUDY",
		"A_MIXTURE_OF_SUN_AND_CLOUDS",
		"HIGH_LEVEL_CLOUDS",
		"MORE_CLOUDS_THAN_SUN",
		"PARTLY_SUNNY",
		"BROKEN_CLOUDS",
		"MOSTLY_CLOUDY",
		"CLOUDY",
		"OVERCAST",
		"LOW_CLOUDS",
		"LIGHT_FOG",
		"FOG",
		"DENSE_FOG",
		"ICE_FOG",
		"SANDSTORM",
		"DUSTSTORM",
		"INCREASING_CLOUDINESS",
		"DECREASING_CLOUDINESS",
		"CLEARING_SKIES",
		"BREAKS_OF_SUN_LATE",
		"EARLY_FOG_FOLLOWED_BY_SUNNY_SKIES",
		"AFTERNOON_CLOUDS",
		"MORNING_CLOUDS",
		"SMOKE",
		"LOW_LEVEL_HAZE",
	}
}
