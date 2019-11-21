#!/bin/bash
set -exo pipefail

telegraf_dir=$(cd $(dirname ${BASH_SOURCE}) && pwd)

function create_security_group() {
  echo "Creating Telegraf scrape security group"
  cf create-security-group telegraf-scrape "${telegraf_dir}/asg.json"
  cf bind-security-group telegraf-scrape system system
}

function download_telegraf() {
  telegraf_version=$(curl -s https://api.github.com/repos/influxdata/telegraf/releases/latest | jq -r .tag_name || "1.12.6")
  telegraf_binary_url="https://dl.influxdata.com/telegraf/releases/telegraf-${telegraf_version}-static_linux_amd64.tar.gz"
  wget -qO- "$telegraf_binary_url" | tar xvz - --strip=2 telegraf/telegraf
}

function create_certificates() {
  mkdir -p certs
  pushd certs > /dev/null
    ca_cert_name=$(credhub find -n metric_scraper_ca --output-json | jq -r .credentials[].name | grep cf)
    credhub generate -n telegraf_scrape_tls -t certificate --ca "$ca_cert_name" -c telegraf_scrape_tls

    credhub get -n telegraf_scrape_tls --output-json | jq -r .value.ca > scrape_ca.crt
    credhub get -n telegraf_scrape_tls --output-json | jq -r .value.certificate > scrape.crt
    credhub get -n telegraf_scrape_tls --output-json | jq -r .value.private_key > scrape.key
  popd > /dev/null
}

function push_telegraf() {
  cf target -o system -s system

  GOOS=linux go build -o confgen
  cf v3-create-app telegraf
  cf v3-apply-manifest -f "${telegraf_dir}/manifest.yml"
  cf v3-push telegraf
}

pushd ${telegraf_dir} > /dev/null
  download_telegraf
  create_security_group
  create_certificates
  push_telegraf
popd > /dev/null