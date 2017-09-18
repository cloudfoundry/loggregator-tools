#!/usr/bin/env bash

# TODO: receive generated id

# JOB_NAME: the name of the job, this should be unique for each dataset
# CYCLES: how many log messages to emit
# DELAY_US: how many microseconds to wait inbetween writing log messages
# CF_SYSTEM_DOMAIN: system domain for communicating with cf api
# CF_USERNAME: cf username
# CF_PASSWORD: cf password
# CF_SPACE: cf space for running test
# CF_ORG: cf org for running test

set -eu

source ./shared.sh

function hammer_url {
    # TODO: add ID to text param
    echo "$(app_url "$(drainspinner_app_name)")?cycles=${CYCLES}&delay=${DELAY_US}us&id=$(test_uuid)"
}

function establish_logs_stream {
    checkpoint "Starting App Logs Stream"

    cf logs "$(drainspinner_app_name)" > output.txt 1>&1 &
    local wait=10
    echo "sleeping for ${wait}s to wait for log stream to be established"
    sleep "$wait"
}

function hammer {
    checkpoint "Writing messages"

    curl "$(hammer_url)" &> /dev/null

    export -f block_until_count_equals_cycles
    if ! timeout "${WAIT:-180}s" bash -ec "block_until_count_equals_cycles"; then
        warning "timed out waiting for all the messages to be received"
    fi
}

function block_until_count_equals_cycles {
    source ./shared.sh
    while true; do
        local count=$(curl -s "$(app_url "$(counter_app_name)")/get/$(test_uuid)")
        if [ "${count:-0}" -ge "$CYCLES" ]; then
            success "received all messages with count $count"
            break
        fi
        echo "waiting to receive all messages, current count $count"
        sleep 5
    done
    exit 0
}

function validate_hammer {
    validate_variables JOB_NAME CYCLES DELAY_US CF_SYSTEM_DOMAIN CF_USERNAME \
        CF_PASSWORD CF_SPACE CF_ORG
}

function main {
    validate_hammer

    checkpoint "Starting Hammer"

    login
    establish_logs_stream
    hammer
}
main
