#!/bin/bash
set -ex

echo "==== updating package index ===="
apt update

echo "==== installing prerequisites ===="
apt install -y apt-transport-https ca-certificates curl software-properties-common gpg

echo "==== Adding Docker's official GPG key and Repo ===="
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | tee /etc/apt/sources.list.d/docker.list
apt update

echo "==== Installing Docker CE ===="
apt install -y docker-ce docker-ce-cli containerd.io

echo "==== Pulling provisioner docker image ===="
docker pull ghcr.io/lachstec/provisioner:dev

echo "==== Starting backend service ===="
docker run -d \
    --restart=always \
    -e DB_HOST="${DB_HOST}" \
    -e DB_PORT="${DB_PORT}" \
    -e DB_USER="${DB_USER}" \
    -e DB_PASS="${DB_PASS}" \
    -e OPENSTACK_IDENTITY_ENDPOINT="${OPENSTACK_IDENTITY_ENDPOINT}" \
    -e OPENSTACK_USER="${OPENSTACK_USER}" \
    -e OPENSTACK_PASS="${OPENSTACK_PASS}" \
    -e OPENSTACK_DOMAIN="${OPENSTACK_DOMAIN}" \
    -e OPENSTACK_TENANT_NAME="${OPENSTACK_TENANT_NAME}" \
    -e CRYPTO_KEY="${CRYPTO_KEY}" \
    -e TRACE_ENDPOINT="${TRACE_ENDPOINT}" \
    -e TRACE_SERVICENAME="${TRACE_SERVICENAME}" \
    -e AUTH0_URL="${AUTH0_URL}" \
    -e AUTH0_AUDIENCE="${AUTH0_AUDIENCE}" \
    -e AUTH0_SECRET="${AUTH0_SECRET}" \
    -p "10000:10000" \
    ghcr.io/lachstec/provisioner:dev

echo "==== Successfully started backend container ===="