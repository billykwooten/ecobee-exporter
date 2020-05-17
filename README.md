# ecobee_exporter

Lots of references from: [https://github.com/dichro/ecobee](https://github.com/dichro/ecobee)

Check him out as well, initial idea from his repository.

## Summary

Ecobee exporter for metrics in Prometheus format

Setting up an ecobee exporter is complicated due to the need to authenticate with the Ecobee API service via tokens.
The first time you run this program it requires some manual steps, however the exporter will subsequently manage its
own authentication afterwards if you give the program a volume somewhere to store passwords and manage it's authorization cache.

## First time running ecobee exporters, read this

1. Create a new app. To do this sign into the ecobee website and look for the `Developer` page, create a new app there. The API key will be used later.
2. Create a volume on your host so we can persist authentication cache
3. Run `docker run -v <volume from step 1>:/db -p 8080:8080 -it billykwooten/ecobee-exporter`
4. Open a browser and go to http://localhost:8080/metrics or `curl -X GET http://localhost:8080/metrics` from another terminal
5. Docker will now print a pin like `Pin is "ig7j"`
6. Go to [https://www.ecobee.com/consumerportal/index.html#/my-apps](https://www.ecobee.com/consumerportal/index.html#/my-apps)
7. Register your app pin from step 4
8. You can now run the container in any way you want, as long as you mount in the volume from step 1.

   Example: [Binary/Docker Run Examples](https://github.com/billykwooten/ecobee_exporter/tree/development#usage)

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
    volumes:
      - /volume1/docker/ecobee-exporter/data:/db
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

## Development

If you'd like to build this yourself you can clone this repo and run:

```
./script/cibuild
```