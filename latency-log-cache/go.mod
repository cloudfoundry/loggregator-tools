module github.com/cloudfoundry-incubator/loggregator-tools/latency

require (
	code.cloudfoundry.org/go-loggregator v7.6.0+incompatible
	code.cloudfoundry.org/log-cache v0.0.0-00010101000000-000000000000
	github.com/cloudfoundry/noaa v2.1.0+incompatible
	github.com/cloudfoundry/sonde-go v0.0.0-20171206171820-b33733203bb4
	github.com/elazarl/goproxy v0.0.0-20181003060214-f58a169a71a5 // indirect
	github.com/gorilla/websocket v1.4.0 // indirect
	github.com/nu7hatch/gouuid v0.0.0-20131221200532-179d4d0c4d8d
)

replace code.cloudfoundry.org/log-cache => code.cloudfoundry.org/log-cache-release/src v2.0.2-0.20190801153242-56147f52a11c+incompatible
