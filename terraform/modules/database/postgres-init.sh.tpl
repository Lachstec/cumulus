#!/bin/bash
set -ex

PG_USER="${pg_user}"
PG_PASSWORD="${pg_password}"
POSTGRES_SUBNET_CIDR="${postgres_subnet_cidr}"

apt update && apt upgrade -y
apt install -y postgresql postgresql-contrib

PG_CONF_DIR=$(find /etc/postgresql -type d -name main | head -n 1)

if [ -z "$PG_CONF_DIR" ]; then
    echo "failed to locate postgres configuration directory"
    exit 1
fi

sed -i "s/#listen_addresses = 'localhost'/listen_addresses = '*'/" "$PG_CONF_DIR/postgresql.conf"
echo "password_encryption = md5" >> "$PG_CONF_DIR/postgresql.conf"
echo "host    all             all             ${postgres_subnet_cidr}            md5" >> "$PG_CONF_DIR/pg_hba.conf"

systemctl restart postgresql
systemctl enable postgresql

sudo -u postgres psql <<EOF
DO \$\$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_roles WHERE rolname = '${pgpool_user}') THEN
      CREATE ROLE ${pgpool_user} WITH LOGIN PASSWORD '${pgpool_password}' CREATEDB REPLICATION SUPERUSER;
   ELSE
      -- In case the role already exists, ensure it has the necessary attributes.
      ALTER ROLE ${pgpool_user} WITH CREATEDB REPLICATION SUPERUSER;
   END IF;
   -- Ensure the password is stored with md5 encryption
   ALTER ROLE ${pgpool_user} WITH PASSWORD '${pgpool_password}';
END
\$\$;
EOF


echo "PostgreSQL setup complete."