# Loggregator Tools

## [HTTPS Drain][https-drain]

The HTTPS Drain is an example HTTPS server that accepts syslog messages
(RFC-5424). This is a useful tool for debugging and monitoring
[cf-syslog-drain-release][cf-syslog-drain-release].
The HTTPS drain can be configured to POST a counter.

## [Syslog Drain][syslog-drain]

The Syslog Drain is an example TCP server that accepts syslog messages
(RFC-5424) via TCP. This is a useful tool for debugging and monitoring
[cf-syslog-drain-release][cf-syslog-drain-release].
The Syslog drain can be configured to POST a counter.

## [Syslog to Datadog][syslog-to-datadog]

The Syslog to Datadog is an example HTTPS server that accepts syslog messages
(RFC-5424) with metrics in the structured data. The metrics will be sent to
datadog.

The Syslog to Datadog application can be configured with `DATADOG_API_KEY` and
`PORT` environment variables.

## [Datadog Accumulator][data-dog-accumulator]

The Datadog Accumulator scrapes Gauge values for a given source-id and send
the results to Datadog.

## [Log Spinner][logspinner]

Log Spinner is a sample CF application that is written in go. It is compatible
with the go-buildpack.

## [JSON Spinner][jsonspinner]

JSON Spinner is a sample CF application that is written in go. It is compatible
with the go-buildpack. It is used by the cf-syslog-drain black box tests.

## [Request Spinner][request-spinner]

Request Spinner reads from Log Cache for load testing purposes. The request
can be configured to hit specific source IDs, at a given cycle and delay.

## [Log Cache Siege][log-cache-siege]

Log Cache Siege is configured with the address of a request-spinner to
instruct request-spinner to hit every available source ID.

## [Slow Consumer][slow-consumer]

The Slow Consumer is a firehose nozzle that will induce the TrafficController
to cut off the nozzle.

## [Post Printer][postprinter]

The post printer is a CF application that prints every request to stderr.

## [Datadog Forwarder][datadog-forwarder]

The Datadog forwarder reads from Log Cache and forwards metrics to Datadog.

## [CF LogMon][cf-logmon]

The CF LogMon performs a blacbox test for measuring message reliability when
running the command cf logs. This is accomplished by writing groups of logs,
measuring the time it took to produce the logs, and then counting the logs
received in the log stream. This is one way to measure message reliability of
the Loggregator system. The results of this test are displayed in a simple UI
and available via JSON and the Firehose.

[https-drain]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/https_drain
[syslog-drain]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog_drain
[cf-syslog-drain-release]: https://github.com/cloudfoundry/cf-syslog-drain-release
[syslog-to-datadog]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog_to_datadog
[data-dog-accumulator]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/experimental/data-dog-accumulator
[logspinner]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/logspinner
[jsonspinner]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/jsonspinner
[request-spinner]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/request-spinner
[log-cache-siege]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/log-cache-siege
[slow-consumer]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/slow_consumer
[postprinter]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/postprinter
[datadog-forwarder]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/log-cache-forwarders/datadog
[cf-logmon]: https://github.com/cloudfoundry-incubator/cf-logmon
