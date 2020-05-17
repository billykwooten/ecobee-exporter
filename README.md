# ecobee_exporter

Lots of references from: [https://github.com/dichro/ecobee](https://github.com/dichro/ecobee)

Check him out as well, initial idea from his repository.

## Summary

Ecobee exporter for metrics in Prometheus format

Setting up an ecobee exporter is complicated due to the need to authenticate with the Ecobee API service via tokens.
The first time you run this program it requires some manual steps, however the exporter will subsequently manage its
own authentication afterwards if you give the program a volume somewhere to store passwords and manage it's authorization cache.

## Configuration

Ecobee exporter can be controlled by both ENV or CLI flags as described below.

| Environment        	       | CLI (`--flag`)              | Default                 	    | Description                                                                                                      |
|----------------------------|-----------------------------|---------------------------- |------------------------------------------------------------------------------------------------------------------|
| `ECOBEE_LISTEN_ADDRESS`           | `listen-address`            | `:9098`                     | The port for /metrics to listen on |
| `ECOBEE_APPKEY`                   | `appkey`                    | `<REQUIRED>`                | Your Application API Key |
| `ECOBEE_CACHEFILE`                     | `cachefile`                      | `/db/auth.cache`              | Cache file to store auth credentials |

## Usage

Binary Usage
```
# Export ecobee metrics from thermostat
./ecobee-exporter --appkey mi4o2n54i0510n4510
```

Docker Usage (recommended method of running)
```
# Export ecobee metrics from thermostat using docker with volume for cache
docker run -d --restart always --name=ecobee-exporter -v /example/persistancedirectory:/db -p 9098:9098 billykwooten/ecobee-exporter --appkey mi4o2n54i0510n4510
```

Docker-compose Usage
```
  ecobee-exporter:
    image: billykwooten/ecobee-exporter
    container_name: ecobee-exporter
    restart: always
    ports:
      - "9098:9098"
    environment:
      - ECOBEE_APPKEY=mi4o2n54i0510n4510
```

Prometheus Scrape Usage
```
scrape_configs:
  - job_name: 'ecobee-exporter'
    scrape_interval: 60s
    static_configs:
      - targets: ['ecobee-exporter:9098']
```
