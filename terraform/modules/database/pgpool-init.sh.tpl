#!/bin/bash
set -ex

PGPOOL_USER="${pgpool_user}"
PGPOOL_PASSWORD="${pgpool_password}"
PGSQL_NODES="${pgsql_nodes}"
POSTGRES_SUBNET_CIDR="${postgres_subnet_cidr}"

# Update and install dependencies
apt update && apt upgrade -y
apt install -y pgpool2

PGPOOL_CONF="/etc/pgpool2/pgpool.conf"
PCP_CONF="/etc/pgpool2/pcp.conf"
POOL_PASSWD="/etc/pgpool2/pool_passwd"

# Ensure pgpool user exists
id -u pgpool &>/dev/null || useradd -r -s /bin/false pgpool

# Set proper ownership for pgpool config directory


# Enable load balancing and master-slave mode
sed -i "s/^#load_balance_mode = off/load_balance_mode = on/" "$PGPOOL_CONF"
sed -i "s/^#master_slave_mode = off/master_slave_mode = on/" "$PGPOOL_CONF"

# Configure PostgreSQL nodes
IFS=',' read -r -a NODES <<< "$${PGSQL_NODES}"
for i in "$${!NODES[@]}"; do
  echo "backend_hostname$${i} = '$${NODES[$${i}]}'" >> "$${PGPOOL_CONF}"
  echo "backend_port$${i} = 5432" >> "$${PGPOOL_CONF}"
  echo "backend_weight$${i} = 1" >> "$${PGPOOL_CONF}"
  echo "backend_data_directory$${i} = '/var/lib/postgresql'" >> "$${PGPOOL_CONF}"
done

# Allow connections from PostgreSQL subnet
echo "host    all             all             $${POSTGRES_SUBNET_CIDR}            md5" >> /etc/pgpool2/pool_hba.conf

# Make Pgpool listen on all interfaces (or you can set specific interfaces)
sed -i "s/^#listen_addresses = 'localhost'/listen_addresses = '*'/g" "$PGPOOL_CONF"

# Configure pgpool user and password
pg_md5 -m -u ${pgpool_user} ${pgpool_password} > "$PCP_CONF"
pg_md5 -m -u ${pgpool_user} ${pgpool_password} > "$POOL_PASSWD"
mkdir -p /var/run/pgpool

# Add failover.sh script
cat << 'EOL' > /etc/pgpool2/failover.sh
#!/bin/bash
set -ex

NODE_ID=$1
FAILED_NODE_HOST=$2
OLD_PRIMARY_ID=$3
NEW_PRIMARY_ID=$4
NEW_PRIMARY_HOST=$5

echo "[FAILOVER] Node $NODE_ID ($FAILED_NODE_HOST) is down. Promoting $NEW_PRIMARY_HOST"

# Promote new primary if it's a standby
if [ "$NODE_ID" -eq "$OLD_PRIMARY_ID" ]; then
  PGPASSWORD="${pgpool_password}" psql -U "${pgpool_user}" -h "$NEW_PRIMARY_HOST" -c "SELECT pg_promote();" || exit 1
fi

exit 0
EOL

# Make the failover.sh script executable
chmod +x /etc/pgpool2/failover.sh

chown -R postgres:postgres /etc/pgpool2
chmod 750 /etc/pgpool2
chmod 640 /etc/pgpool2/pgpool.conf
chmod 600 /etc/pgpool2/pcp.conf /etc/pgpool2/pool_passwd
chmod 750 /etc/pgpool2/failover.sh

# Update pgpool.conf with failover command and health check parameters
echo "failover_command = '/etc/pgpool2/failover.sh %d %h %P %M %H'" >> /etc/pgpool2/pgpool.conf
echo "health_check_period = 5" >> /etc/pgpool2/pgpool.conf
echo "health_check_timeout = 10" >> /etc/pgpool2/pgpool.conf
echo "health_check_user = '${pgpool_user}'" >> /etc/pgpool2/pgpool.conf
echo "health_check_password = '${pgpool_password}'" >> /etc/pgpool2/pgpool.conf

# Restart pgpool2 to apply changes
systemctl restart pgpool2
systemctl enable pgpool2

echo "pgpool2 setup completed."