## CF-pushable Prometheus server

The `./push.sh` downloads the latest prometheus version,sets up the 
appropriate security groups to allow network access, configures 
certificates using credhub for and finaly pushes prometheus server. 

This repo includes an example `prometheus.yml`, and the IPs under `static_configs` 
should be replaced by IPs from `bosh vms` of your deployment. The information
for the security group is stored in the `asg.json` file. 

After pushing the app, visit the Prometheus dashboard at
`prometheus.<system_domain>/graph`.
