#! /usr/bin/env bash

set -ex

function create_app {
  set -e
  local appname=$1
  local droplet_guid=$2

  cf v3-create-app ${appname}
  new_droplet_guid=$(cf curl /v3/droplets?source_guid=${droplet_guid} -d '{
  "relationships": {
    "app": {
      "data": {
        "guid": "'$(cf app ${appname} --guid)'"
      }
    }
  }' | jq -r .guid)

  wait_for_droplet ${appname} ${new_droplet_guid}

  cf v3-set-droplet ${appname} -d ${new_droplet_guid} 2>/dev/null

  process_guid=$(cf curl /v3/apps/$(cf app ${appname} --guid)/processes | jq -r .resources[0].guid)
  cf curl /v3/processes/${process_guid} -X PATCH -d '{"command": "./continuous_log_emitter"}'
  cf curl /v3/processes/${process_guid}/actions/scale -X POST -d '{"memory_in_mb": 32, "disk_in_mb": 32}'
  configure_and_start ${appname}
}

function configure_and_start {
    appname=$1

    cf set-env ${appname} EMIT_INTERVAL $emit_interval
    cf start ${appname}
}

function wait_for_droplet {
  local appname=$1
  local droplet_guid=$2

  set +x
    while ! cf v3-droplets "${appname}" 2>/dev/null | grep "${droplet_guid}" | grep staged; do
      sleep .1s
    done
  set -x
}

if [ "$#" -ne 4 ]; then
    echo "usage: $0 <app name> <drain ip> <emit interval> <droplet guid>"
    exit 1
fi

name=$1
drain_dest_ip=$2
emit_interval=$3
droplet_guid=$4

create_app $name $droplet_guid

emit_interval=2ms

if [ -n "$EMIT_INTERVAL" ]; then
	emit_interval=$EMIT_INTERVAL
fi

for drain in $(cf drains | grep $name | awk '{print $2}'); do
    cf delete-drain --force $drain
done

if [ -z "$DRAIN_COUNT" ]; then
    cf drain $name \
        "https$SCHEME_SUFFIX://$drain_dest_ip:8080?drain_num=1" \
        --type logs
else
    for i in $(seq $DRAIN_COUNT); do
        cf drain $name \
            "https$SCHEME_SUFFIX://$drain_dest_ip:8080?drain_num=$i" \
            --type logs
    done
fi
