# Loggregator Tools

## [HTTPS Drain](https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/https_drain)

The HTTPS Drain is an example HTTPS server that accepts syslog messages
(RFC-5424). This is a useful tool for debugging and monitoring
[cf-syslog-drain-release](https://github.com/cloudfoundry/cf-syslog-drain-release).
The HTTPS drain can be configured to POST a counter.

## [Syslog Drain](https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog_drain)

The Syslog Drain is an example TCP server that accepts syslog messages
(RFC-5424) via TCP. This is a useful tool for debugging and monitoring
[cf-syslog-drain-release](https://github.com/cloudfoundry/cf-syslog-drain-release).
The Syslog drain can be configured to POST a counter.

## [Syslog to Datadog](https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog_to_datadog)

The Syslog to Datadog is an example HTTPS server that accepts syslog messages
(RFC-5424) with metrics in the structured data. The metrics will be sent to
datadog.

The Syslog to Datadog application can be configured with `DATADOG_API_KEY` and
`PORT` environment variables.

## [Log Spinner](https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/logspinner)

Log Spinner is a sample CF application that is written in go. It is compatible
with the go-buildpack.

## [Slow Consumer](https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/slow_consumer)

The Slow Consumer is a firehose nozzle that will induce the TrafficController
to cut off the nozzle.

## [Post Printer](https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/postprinter)

The post printer is a CF application that prints every request to stderr.

## [CF LogMon](https://github.com/cloudfoundry-incubator/cf-logmon)

The CF LogMon performs a blacbox test for measuring message reliability when
running the command cf logs. This is accomplished by writing groups of logs,
measuring the time it took to produce the logs, and then counting the logs
received in the log stream. This is one way to measure message reliability of
the Loggregator system. The results of this test are displayed in a simple UI
and available via JSON and the Firehose.
