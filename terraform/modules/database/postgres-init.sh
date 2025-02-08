#!/bin/bash
set -ex

apt update && apt upgrade -y
apt install -y postgresql postgresql-contrib

PG_CONF_DIR=$(find /etc/postgresql -type d -name main | head -n 1)

if [ -z "$PG_CONF_DIR" ]; then
    echo "failed to locate postgres configuration directory"
    exit 1
fi

sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" "$PG_CONF_DIR/postgresql.conf"
echo "host    all             all             ${POSTGRES_SUBNET_CIDR}            md5" >> "$PG_CONF_DIR/pg_hba.conf"

systemctl restart postgresql
systemctl enable postgresql
systemctl status postgresql --no-pager

