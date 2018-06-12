### Syslog Nozzle

This is a syslog nozzle that reads from the RLP and sends the envelopes to a
syslog destination.

It filters and routes the logs via a `namespace` tag in the envelope to the
apporpriate destination.

__Hint:__ This is meant to be deployed in Kubernetes via [this docker
image][syslog-nozzle-docker-image] and [this kubernetes
deployment][syslog-nozzle-deployment].



[syslog-nozzle-docker-image]: https://github.com/cloudfoundry/loggregator-ci/blob/master/docker-images/syslog-nozzle/Dockerfile
[syslog-nozzle-deployment]: https://github.com/cloudfoundry-incubator/loggregator-k8s-deployment/blob/master/optional/deployments/syslog-nozzle.yml
