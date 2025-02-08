#!/bin/bash
set -ex

NODE_ID=$1
FAILED_NODE_HOST=$2
OLD_PRIMARY_ID=$3
NEW_PRIMARY_ID=$4
NEW_PRIMARY_HOST=$5

echo "[FAILOVER] Node $NODE_ID ($FAILED_NODE_HOST) is down. Promoting $NEW_PRIMARY_HOST" >> /var/log/pgpool/failover.log

# Promote new primary if it's a standby
if [ "$NODE_ID" -eq "$OLD_PRIMARY_ID" ]; then
  echo "[FAILOVER] Promoting new primary: $NEW_PRIMARY_HOST" >> /var/log/pgpool/failover.log
  PGPASSWORD="${PGPOOL_PASSWORD}" psql -U "${PGPOOL_USER}" -h "$NEW_PRIMARY_HOST" -c "SELECT pg_promote();" || exit 1
fi

exit 0