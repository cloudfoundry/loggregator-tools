
function login {
    checkpoint "Logging into CF"

    if [ "${login_has_occurred:-}" = "" ]; then
        cf login \
            -a api."$CF_SYSTEM_DOMAIN" \
            -u "$CF_USERNAME" \
            -p "$CF_PASSWORD" \
            -s "$CF_SPACE" \
            -o "$CF_ORG" \
            --skip-ssl-validation # TODO: consider passing this in as a param
    fi
    export login_has_occurred=true
}

function validate_variables {
    for var in "$@"; do
        local value=${!var:-}
        case "$var" in
            DRAIN_VERSION)
                if [ "$value" != "1.0" ] && [ "$value" != "2.0" ]; then
                    error "$var must be either \"1.0\" or \"2.0\""
                    return 1
                fi
                ;;
            DRAIN_TYPE)
                if [ "$value" != "syslog" ] && [ "$value" != "https" ]; then
                    error "$var must be either \"syslog\" or \"https\""
                    return 1
                fi
                ;;
            SINK_DEPLOY)
                if [ "$value" = "standalone"  ]; then
                    for v in "EXTERNAL_DRAIN_HOST EXTERNAL_DRAIN_PORT EXTERNAL_COUNTER_PORT"; do
                        if [ "$v" = "" ]; then
                            error "$v must be set for a standalone $var"
                            return 1
                        fi
                    done
                elif [ "$value" != "cf" ] && [ "$value" != "standalone"  ]; then
                    error "$var must be either \"cf\" or \"standalone\""
                    return 1
                elif [ "$value" = "standalone" ] && [ "$DRAIN_TYPE" = "https" ]; then
                    error "$var as standalone is not supported for HTTPS
                    drains"
                    return 2
                fi
                ;;
            CYCLES)
                if [ ! "$value" -gt 0 ]; then
                    error "$var must be a positive number"
                    return 1
                fi
                ;;
            DELAY_US)
                if [ ! "$value" -ge 0 ]; then
                    error "$var must be a nonnegative number"
                    return 1
                fi
                ;;
            *)
                if [ "$value" = "" ]; then
                    error "$var needs to be set"
                    return 1
                fi
                ;;
        esac
    done
}

function checkpoint {
    echo
    echo -e "\e[95m##### $1 #####\e[0m"
}

function error {
    echo -e "\e[91m$1\e[0m"
}

function warning {
    echo -e "\e[93m$1\e[0m"
}

function success {
    echo -e "\e[92m$1\e[0m"
}

function is_standalone {
    [ "$SINK_DEPLOY" = "standalone" ]
}

function app_url {
    local app_name="$1"

    if [ "$app_name" = "" ]; then
        echo app name not provided
        exit 22
    fi

    local guid=$(cf app "$app_name" --guid)
    local route_data=$(cf curl "/v2/apps/$guid/routes")
    local domain_url=$(echo "$route_data" | jq .resources[0].entity.domain_url --raw-output)
    local domain_name=$(cf curl "$domain_url" | jq .entity.name --raw-output)

    local port=$(echo "$route_data" | jq .resources[0].entity.port --raw-output)
    if [ "$port" != "null" ]; then
        # this app uses tcp routing
        echo "$domain_name:$port"
    else
        local host=$(echo "$route_data" | jq .resources[0].entity.host --raw-output)
        echo "$host.$domain_name"
    fi
}

function get_prime_count {
    if $(is_standalone); then
        echo "$EXTERNAL_DRAIN_HOST:$EXTERNAL_COUNTER_PORT/get-prime/$(test_uuid)" && return
    fi

    echo "$(app_url "$(counter_app_name)")/get-prime/$(test_uuid)"
}

function counter_endpoint {
    if $(is_standalone); then
        echo "$EXTERNAL_DRAIN_HOST:$EXTERNAL_COUNTER_PORT/get/$(test_uuid)" && return
    fi

    echo "$(app_url "$(counter_app_name)")/get/$(test_uuid)"
}

function drain_app_name {
    echo "ss-smoke-drain-$JOB_NAME"
}

function drainspinner_app_name {
    echo "ss-smoke-drainspinner-$JOB_NAME"
}

function counter_app_name {
    echo "ss-smoke-counter-$JOB_NAME"
}

function syslog_drain_service_name {
    echo "ss-smoke-${JOB_NAME}-${DRAIN_VERSION}"
}

function syslog_drain_service_url {
    if $(is_standalone); then
        echo "$DRAIN_TYPE://$EXTERNAL_DRAIN_HOST:$EXTERNAL_DRAIN_PORT" && return
    fi

    echo "$DRAIN_TYPE://$(app_url $(drain_app_name))/drain"
}

function test_uuid {
    if [ ! -e /tmp/test_uuid_${CYCLES}_${DELAY_US} ]; then
        cat /proc/sys/kernel/random/uuid > "/tmp/test_uuid_${CYCLES}_${DELAY_US}"
    fi

    cat "/tmp/test_uuid_${CYCLES}_${DELAY_US}"
}
