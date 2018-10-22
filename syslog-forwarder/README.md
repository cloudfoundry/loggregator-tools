# Syslog Forwarder App

The syslog forwarder app creates a connection to the Loggregator Reverse Log
Proxy (RLP) for each app in the space it's deployed to

## Deploy The Forwarder

First, you'll need to make an app manifest so you can `cf push`.
Your best bet is to edit the `scripts/manifest.yml` file.
The only thing you need to configure are the following environment variables:

```
  SOURCE_HOSTNAME: <The hostname that will be reported in the output syslog bodies>
  INCLUDE_SERVICES: <Whether to include on-demand-service instance in the deployed space>
  SYSLOG_URL: <The endpoint to send your logs to. https, syslog, and syslog-tls are supported>
```

From the `syslog-forwarder` directory in this repository, run:

```
./scripts/deploy.sh
```

The script will build and push the syslog forwarder as an app called `syslog-forwarder` in
whichever space you happen to have targeted.

### If you want to build the forwarder without deploying it
From the `syslog-forwarder` directory in this repository, run:

```
./scripts/build-forwarder.sh
```

build-forwarder results in a zip file containing two binaries that make up
the syslog-forwarder with a `run.sh` script to orchestrate their execution.

To push the zipped forwarder run:

```
cf push <app-name> -b binary_buildpack -c ./run.sh -u process -p forwarder.zip
```
