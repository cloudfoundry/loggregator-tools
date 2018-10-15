#! /bin/bash

git_root=$(git rev-parse --show-toplevel)

pushd $git_root/syslog-forwarder/scripts
    ./build-forwarder.sh
    cf push syslog-forwarder --no-start -f manifest.yml
    cf set-env syslog-forwarder REFRESH_TOKEN "$(cat ~/.cf/config.json | jq -r .RefreshToken)"
    cf set-env syslog-forwarder CLIENT_ID "$(cat ~/.cf/config.json | jq -r .UAAOAuthClient)"
    cf set-env syslog-forwarder SKIP_SSL_VALIDATION "$(cat ~/.cf/config.json | jq -r .SSLDisabled)"
    cf start syslog-forwarder
popd
