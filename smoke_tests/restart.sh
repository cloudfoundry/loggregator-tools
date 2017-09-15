#!/usr/bin/env bash
set -eu

source ./shared.sh

function main {
    cf restart "$(drainspinner_app_name)"
    cf restart "$(drain_app_name)"
}
main
