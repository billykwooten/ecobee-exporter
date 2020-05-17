# ecobee_exporter

Lots of references from: [https://github.com/dichro/ecobee](https://github.com/dichro/ecobee)

Check him out as well, initial idea from his repository.

## Summary

Ecobee exporter for metrics in Prometheus format

Setting up an ecobee exporter is complicated due to the need to authenticate with the Ecobee API service via tokens.
The first time you run this program it requires some manual steps, however the exporter will subsequently manage its
own authentication afterwards if you give the program a volume somewhere to store passwords and manage it's authorization cache.

