## CF-pushable Prometheus server

The `./push.sh` sets up the appropriate security groups to allow scraping for
Prometheus as well as pushing the Prometheus server. This repo includes an
example `prometheus.yml`, and the IPs under `static_configs` should be
replaced by IPs from `bosh vms` in your deployment.

After pushing the app, visit the Prometheus dashboard at
`prometheus.<system_domain>/graph`.
