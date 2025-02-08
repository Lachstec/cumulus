#!/bin/bash
set -ex

PGPOOL_USER="${pgpool_user}"
PGPOOL_PASSWORD="${pgpool_password}"
PGSQL_NODES="${pgsql_nodes}"
PGPOOL_VIP="${pgpool_vip}"
POSTGRES_SUBNET_CIDR="${postgres_subnet_cidr}"

apt update && apt upgrade -y
apt install -y pgpool2

PGPOOL_CONF="/etc/pgpool2/pgpool.conf"
PCP_CONF="/etc/pgpool2/pcp.conf"
POOL_PASSWD="/etc/pgpool2/pool_passwd"

sed -i "s/^#load_balance_mode = off/load_balance_mode = on/" "$PGPOOL_CONF"
sed -i "s/^#master_slave_mode = off/master_slave_mode = on/" "$PGPOOL_CONF"

IFS=',' read -r -a NODES <<< "$PGSQL_NODES"
for i in "${!NODES[@]}"; do
  echo "backend_hostname$i = '${NODES[$i]}'" >> "$PGPOOL_CONF"
  echo "backend_port$i = 5432" >> "$PGPOOL_CONF"
  echo "backend_weight$i = 1" >> "$PGPOOL_CONF"
  echo "backend_data_directory$i = '/var/lib/postgresql'" >> "$PGPOOL_CONF"
done

echo "host    all             all             ${POSTGRES_SUBNET_CIDR}            md5" >> /etc/pgpool2/pool_hba.conf

echo "${PGPOOL_USER}:$(pg_md5 ${PGPOOL_PASSWORD})" > "$PCP_CONF"
chmod 600 "$PCP_CONF"

echo "${PGPOOL_USER}:$(pg_md5 ${PGPOOL_PASSWORD})" > "$POOL_PASSWD"
chmod 600 "$POOL_PASSWD"

systemctl restart pgpool2
systemctl enable pgpool2
echo "pgpool2 setup completed."