BEGIN;

ALTER TABLE mch_provisioner.servers
DROP CONSTRAINT fk_ssh_key;

ALTER TABLE mch_provisioner.servers
DROP COLUMN ssh_key;

ALTER TABLE mch_provisioner.servers
    ADD COLUMN ssh_key BYTEA;

DROP TABLE mch_provisioner.keypairs;

COMMIT;