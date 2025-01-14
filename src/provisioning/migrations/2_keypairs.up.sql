BEGIN;

CREATE TABLE mch_provisioner.keypairs(
    name VARCHAR PRIMARY KEY,
    public_key BYTEA,
    private_key BYTEA
);

ALTER TABLE mch_provisioner.servers
DROP COLUMN ssh_key;

ALTER TABLE mch_provisioner.servers
ADD COLUMN ssh_key VARCHAR;

ALTER TABLE mch_provisioner.servers
ADD CONSTRAINT fk_ssh_key
FOREIGN KEY (ssh_key)
REFERENCES mch_provisioner.keypairs(name);

COMMIT;