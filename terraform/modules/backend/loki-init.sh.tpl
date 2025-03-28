#!/bin/bash
set -ex

echo "==== updating package index ===="
apt update

echo "==== installing prerequisites ===="
apt install -y apt-transport-https ca-certificates curl gnupg-agent software-properties-common

echo "==== Adding Docker GPG key and repository ===="
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list
apt update

echo "==== Installing Docker CE ===="
apt install -y docker-ce docker-ce-cli containerd.io

echo "==== Creating a dedicated Docker network ===="
docker network create loki-net || true

cat <<EOL > /tmp/loki-config.yaml
auth_enabled: false

server:
  http_listen_port: 3100

ingester:
  lifecycler:
    ring:
      kvstore:
        store: inmemory
      replication_factor: 1
  final_sleep: 0s

schema_config:
  configs:
    - from: 2020-10-15
      store: boltdb-shipper
      object_store: swift
      schema: v11
      index:
        prefix: index_
        period: 24h

storage_config:
  boltdb_shipper:
    active_index_directory: /loki/index
    cache_location: /loki/cache
    shared_store: swift
    chunk_target_size: 1048576
  swift:
    auth_url: "${openstack_auth_url}"
    username: "${openstack_username}"
    password: "${openstack_password}"
    tenant_name: "${openstack_tenant_name}"
    domain_name: "${openstack_domain_name}"
    region_name: "${openstack_region_name}"
    container_name: "${loki_container_name}"
    insecure_skip_verify: true  # This will skip SSL certificate verification

limits_config:
  enforce_metric_name: false
  reject_old_samples: true
  reject_old_samples_max_age: 168h
  allow_structured_metadata: true
EOL

echo "==== Starting Loki Container ===="
docker run -d --name loki --network loki-net -p 3100:3100 grafana/loki:3.3.2 \
  -config.expand-env=true \
  -config.file=/etc/loki/local-config.yaml || true

echo "==== Creating Grafana datasource file ===="
cat <<EOL > /tmp/datasource.yaml
apiVersion: 1
datasources:
  - name: Loki
    type: loki
    access: proxy
    url: http://loki:3100
    isDefault: true
    editable: true
EOL

echo "==== Starting Grafana Container ===="
docker run -d --name grafana --network loki-net -p 3000:3000 \
  -v /tmp/datasource.yaml:/etc/grafana/provisioning/datasources/datasource.yaml \
  grafana/grafana:latest