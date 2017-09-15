#!/usr/bin/env bash

# JOB_NAME: the name of the job, this should be unique for each dataset
# DRAIN_VERSION: 1.0 or 2.0
# CF_SYSTEM_DOMAIN: system domain for communicating with cf api
# CF_USERNAME: cf username
# CF_PASSWORD: cf password
# CF_SPACE: cf space for running test
# CF_ORG: cf org for running test

set -eu

source ./shared.sh

function validate_teardown {
    validate_variables JOB_NAME DRAIN_VERSION CF_SYSTEM_DOMAIN CF_USERNAME \
        CF_PASSWORD CF_SPACE CF_ORG
}

function main {
    validate_teardown

    checkpoint "Tearing Down Apps and Services"

    login
    cf delete "$(drainspinner_app_name)" -r -f
    cf delete "$(drain_app_name)" -r -f
    cf delete "$(counter_app_name)" -r -f
    cf delete-service "$(syslog_drain_service_name)" -f
}
main
