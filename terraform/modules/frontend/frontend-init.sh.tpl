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

echo "==== Pulling frontend docker image ===="
docker pull ghcr.io/lachstec/frontend:dev

echo "==== Starting frontend service ===="
docker run -d \
    --restart=always \
    -e PUBLIC_AUTH_DOMAIN="${PUBLIC_AUTH_DOMAIN}" \
    -e PUBLIC_AUTH_CLIENT_ID="${PUBLIC_AUTH_CLIENT_ID}" \
    -e PUBLIC_AUTH_AUDIENCE="${PUBLIC_AUTH_AUDIENCE}" \
    -e PUBLIC_AUTH_CACHE_LOCATION="${PUBLIC_AUTH_CACHE_LOCATION}" \
    -e PUBLIC_BACKEND_URL="${PUBLIC_BACKEND_URL}" \
    -e PUBLIC_REQUESTER_NAME="${PUBLIC_REQUESTER_NAME}" \
    -p "80:3000" \
    ghcr.io/lachstec/frontend:dev

echo "==== Successfully started frontend container ===="