#!/usr/bin/env bash

# TODO: receive generated id

# JOB_NAME: the name of the job, this should be unique for each dataset
# CYCLES: how many log messages to emit
# DELAY_US: how many microseconds to wait inbetween writing log messages
# DATADOG_API_KEY:
# DRAIN_TYPE: syslog or https
# DRAIN_VERSION: 1.0 or 2.0
# CF_SYSTEM_DOMAIN: system domain for communicating with cf api
# CF_USERNAME: cf username
# CF_PASSWORD: cf password
# CF_SPACE: cf space for running test
# CF_ORG: cf org for running test
# CF_APP_DOMAIN: tcp or https app domain based on DRAIN_TYPE

set -u

source ./shared.sh

function kill_cf {
    pkill cf || true
}

function datadog_url {
    echo "https://app.datadoghq.com/api/v1/series?api_key=$1"
}

function post_to_datadog {
    echo "posting to datadog. name = $1, timestamp = $2, value = $3"
    local payload
    payload=$(cat <<JSON
{
    "series": [{
        "metric": "smoke_test.ss.loggregator.$1",
        "points": [[$2, $3]],
        "type": "gauge",
        "host": "$CF_SYSTEM_DOMAIN",
        "tags": [
            "drain_version:$DRAIN_VERSION",
            "drain_type:$DRAIN_TYPE",
            "job_name:$JOB_NAME",
            "cycles:$CYCLES",
            "delay_us:$DELAY_US"
        ]
    }]
}
JSON
)

for api_key in $DATADOG_API_KEY; do
    curl -X POST -H "Content-type: application/json" -d "$payload" "$(datadog_url "$api_key")"
        echo
done
}

function validate_report {
    validate_variables JOB_NAME CYCLES DELAY_US DATADOG_API_KEY DRAIN_TYPE \
        DRAIN_VERSION CF_SYSTEM_DOMAIN CF_USERNAME CF_PASSWORD CF_SPACE \
        CF_ORG CF_APP_DOMAIN
}

function main {
    validate_report

    checkpoint "Reporting Results"

    kill_cf

    local msg_count
    if [ -e output.txt ]; then
        msg_count=$(grep -c -E 'APP.+msgCount' output.txt)
    else
        error "output.txt was not created"
    fi

    local drain_msg_count
    drain_msg_count=$(curl -s "$(app_url "$(counter_app_name)")/get/$(test_uuid)")

    currenttime=$(date +%s)

    post_to_datadog "msg_count" "$currenttime" "$msg_count"
    post_to_datadog "drain_msg_count" "$currenttime" "$drain_msg_count"
    post_to_datadog "delay" "$currenttime" "$DELAY_US"
    post_to_datadog "cycles" "$currenttime" "$CYCLES"

    if [ "$msg_count" -gt "$CYCLES" ]; then
        error "msg_count ($msg_count) was > CYCLES ($CYCLES)"
        exit 1
    fi

    if [ "$drain_msg_count" -gt "$CYCLES" ]; then
        error "drain_msg_count ($drain_msg_count) was > CYCLES ($CYCLES)"
        exit 1
    fi

    if [ "$msg_count" -eq 0 ]; then
        error "message count was zero, sad"
        exit 1
    fi
    if [ "$drain_msg_count" -eq 0 ]; then
        error "drain count was zero, sad"
        exit 1
    fi
}
main
