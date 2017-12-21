#!/usr/bin/env bash

# JOB_NAME: the name of the job, this should be unique for each dataset
# DRAIN_TYPE: syslog or https
# DRAIN_VERSION: 1.0 or 2.0
# CF_SYSTEM_DOMAIN: system domain for communicating with cf api
# CF_USERNAME: cf username
# CF_PASSWORD: cf password
# CF_SPACE: cf space for running test
# CF_ORG: cf org for running test
# CF_APP_DOMAIN: tcp or https app domain based on DRAIN_TYPE
# SINK_DEPLOY: syslogs drain into a cf pushed app or our standalone TCP Server

set -eu

source ./shared.sh

function ensure_counter_app {
    checkpoint "Ensuring Counter App is Pushed"

    if ! cf app "$(counter_app_name)" &> /dev/null; then
        push_counter_app
    else
        restart_counter_app
    fi
}

function ensure_drain_app {
    checkpoint "Ensuring Drain App is Pushed"

    if ! cf app "$(drain_app_name)" &> /dev/null; then
        push_drain_app
    else
        restart_drain_app
    fi
}

function ensure_spinner_app {
    checkpoint "Ensuring Spinner App is Pushed"

    if ! cf app "$(drainspinner_app_name)" &> /dev/null; then
        push_spinner_app
    else
        restart_spinner_app
    fi
}

function restart_counter_app {
    checkpoint "Restarting Counter App"

    cf restart "$(counter_app_name)"
}

function restart_drain_app {
    checkpoint "Restarting Drain App"

    cf restart "$(drain_app_name)"
}

function restart_spinner_app {
    checkpoint "Restarting Spinner App"

    cf restart "$(drainspinner_app_name)"
}

function push_counter_app {
    checkpoint "Pushing Counter App"

    pushd ../counter
        if ! [ -e ./counter ]; then
            GOOS=linux go build
        fi
        cf push "$(counter_app_name)" -c ./counter -b binary_buildpack -m 128M
    popd
}

function push_drain_app {
    checkpoint "Pushing Drain App"

    pushd "../${DRAIN_TYPE}_drain"
        if ! [ -e "./${DRAIN_TYPE}_drain" ]; then
            GOOS=linux go build
        fi
        cf push "$(drain_app_name)" \
            -c "./${DRAIN_TYPE}_drain" \
            -b binary_buildpack \
            --no-route \
            --no-start \
            -m 128M \
            --health-check-type none
        cf set-env "$(drain_app_name)" \
            COUNTER_URL \
            "http://$(app_url "$(counter_app_name)")"

        if [ "$DRAIN_TYPE" = "syslog" ]; then
            cf map-route "$(drain_app_name)" "$CF_APP_DOMAIN" --random-port
        else
            cf map-route "$(drain_app_name)" "$CF_APP_DOMAIN" --hostname "$(drain_app_name)"
        fi

        cf start "$(drain_app_name)"
    popd
}

function push_spinner_app {
    checkpoint "Pushing Spinner App"

    pushd ../jsonspinner
        if ! [ -e ./jsonspinner ]; then
            GOOS=linux go build
        fi
        for i in {1..5}; do
            if cf push "$(drainspinner_app_name)" -c ./jsonspinner -b binary_buildpack -m 128M; then
                break
            fi
            sleep 5
        done
    popd
}

function ensure_drain_service {
    checkpoint "Ensuring Drain Service Exists"

    if ! cf service "$(syslog_drain_service_name)" &> /dev/null; then
        create_drain_service
    fi
}

function create_drain_service {
    checkpoint "Creating Drain Service at $(syslog_drain_service_url)"

    cf create-user-provided-service \
        "$(syslog_drain_service_name)" \
        -l "$(syslog_drain_service_url)"
}

function bind_service {
    checkpoint "Binding Service to Drainspinner App"

    cf bind-service "$(drainspinner_app_name)" "$(syslog_drain_service_name)"
}

function prime_service_binding {
    checkpoint "Priming Service Binding"

    # start emitting prime messages
    local timeout_seconds=300
    curl "$(app_url "$(drainspinner_app_name)")?cycles=$timeout_seconds&delay=1s&primer=true&id=$(test_uuid)" &> /dev/null

    export -f block_until_count
    if ! timeout "${timeout_seconds}s" bash -ec "block_until_count"; then
        error "unable to prime the syslog drain binding"
        exit 5
    fi
}

function block_until_count {
    # wait for first prime message and exit
    source ./shared.sh
    while true; do
        local count=$(curl -s "$(get_prime_count)")
        if [ "${count:-0}" -gt 0 ]; then
            success "received primer message, binding has been setup"
            break
        fi
        echo waiting for primer message, binding is not setup
        sleep 5
    done
    exit 0
}

function validate_push {
    validate_variables JOB_NAME DRAIN_TYPE DRAIN_VERSION CF_SYSTEM_DOMAIN \
        SINK_DEPLOY CF_USERNAME CF_PASSWORD CF_SPACE CF_ORG CF_APP_DOMAIN
}

function main {
    validate_push

    checkpoint "Starting Push"

    login
    if $(! is_standalone); then
        ensure_counter_app
        ensure_drain_app
    fi

    ensure_spinner_app
    ensure_drain_service
    bind_service

    prime_service_binding
}
main
