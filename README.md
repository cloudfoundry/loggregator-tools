# Loggregator Tools

## [Big Logger][biglogger]

Sends big logs. Allows user to pass in size, frequency, and receiver port as
arguments.

## [CF LogMon][cf-logmon]

The CF LogMon performs a blacbox test for measuring message reliability when
running the command cf logs. This is accomplished by writing groups of logs,
measuring the time it took to produce the logs, and then counting the logs
received in the log stream. This is one way to measure message reliability of
the Loggregator system. The results of this test are displayed in a simple UI
and available via JSON and the Firehose.

## [Counter][counter]

Web-based counter API used by the [Smoke Tests](#smoke-tests).

## [Datadog Accumulator][data-dog-accumulator]

The Datadog Accumulator scrapes Gauge values for a given source-id and send
the results to Datadog.

## [Doppler Client][dopplerclient]

The Doppler Client will connect to the V2 egress API of the Doppler and read
all the logs/metrics from that single Doppler.

The Doppler Client application must be configured with `DATADOG_API_KEY` and
`DOPPLER_ADDR`, `SHARD_ID`, `CA_PATH`, `CERT_PATH`, `KEY_PATH` environment
variables.

## [Dummy Metron][dummymetron]

A NOOP metron (loggregator agent).

## [Echo][echo]

Contains an HTTP and TCP echo server. Both of which will simply print whatever
they receive.

## [Emitter][emitter]

Emits a V2 Counter Envelope to the Loggregator Agent at a one second interval.

## [Envelope Emitter][envelopeemitter]

Emits a V2 Gauge Envelope to the Loggregator Agent at a one second interval.

## [HTTPS Drain][https-drain]

The HTTPS Drain is an example HTTPS server that accepts syslog messages
(RFC-5424). This is a useful tool for debugging and monitoring
[cf-syslog-drain-release][cf-syslog-drain-release].
The HTTPS drain can be configured to POST a counter.

## [JSON Spinner][jsonspinner]

JSON Spinner is a sample CF application that is written in go. It is compatible
with the go-buildpack. It is used by the cf-syslog-drain black box tests.

## [Latency][latency]

Measures the latency from when an application running on CF emits a log to the
time that log is egressed via the V1 app stream.

## [Log Spinner][logspinner]

Log Spinner is a sample CF application that is written in go. It is compatible
with the go-buildpack.

## [Post Printer][postprinter]

The post printer is a CF application that prints every request to stderr.

## [Reliability][reliability]

Reliability is a tool for measuring the Loggregator Firehose. It consists of
two parts, a server and a worker. You should deploy the same number of workers
as you have Traffic Controllers. The server will communicate with the workers
via WebSockets.

## [RLP Reader][rlpreader]

The RLP Reader connects to and reads from the Loggregator Reverse Log Proxy.
It can be ran with a delay, counter name filter, gauge name filter,
preferredTags, deterministicName and/or shardID.

## [RLP Type Reader][rlptypereader]

The RLP Reader connects to and reads from the Loggregator Reverse Log Proxy.
It can be ran with to only receive one or more different types of Envelopes.

## [Slow Consumer][slow-consumer]

The Slow Consumer is a firehose nozzle that will induce the TrafficController
to cut off the nozzle.

## [Smoke Tests][smoke-tests]

Smoke Tests for CF Syslog Drain release which measures reliability at different
rates.

## [Syslog Drain][syslog-drain]

The Syslog Drain is an example TCP server that accepts syslog messages
(RFC-5424) via TCP. This is a useful tool for debugging and monitoring
[cf-syslog-drain-release][cf-syslog-drain-release].
The Syslog drain can be configured to POST a counter.

## [Syslog Forwarder][syslog-forwarder]

Reads logs from the Loggregator V2 API via the RLP (Reverse Log Proxy) Gateway
and writes logs for a configured source ID to a syslog endpoint. The syslog
drains can be configured to use TCP, TCP w/ TLS or HTTPS.

## [Syslog to Datadog][syslog-to-datadog]

The Syslog to Datadog is an example HTTPS server that accepts syslog messages
(RFC-5424) with metrics in the structured data. The metrics will be sent to
datadog.

The Syslog to Datadog application can be configured with `DATADOG_API_KEY` and
`PORT` environment variables.

[biglogger]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/biglogger
[cf-logmon]: https://github.com/cloudfoundry-incubator/cf-logmon
[cf-syslog-drain-release]: https://github.com/cloudfoundry/cf-syslog-drain-release
[counter]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/counter
[data-dog-accumulator]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/experimental/data-dog-accumulator
[dopplerclient]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/dopplerclient
[dummymetron]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/dummymetron
[echo]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/echo
[emitter]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/emitter
[envelopeemitter]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/envelopeemitter
[https-drain]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/https_drain
[jsonspinner]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/jsonspinner
[latency]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/latency
[logspinner]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/logspinner
[postprinter]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/postprinter
[reliability]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/reliability
[rlpreader]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/rlpreader
[rlptypereader]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/rlptypereader
[slow-consumer]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/slow_consumer
[smoke-tests]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/smoke_tests
[syslog-drain]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog_drain
[syslog-forwarder]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog-forwarder
[syslog-to-datadog]: https://github.com/cloudfoundry-incubator/loggregator-tools/tree/master/syslog_to_datadog
