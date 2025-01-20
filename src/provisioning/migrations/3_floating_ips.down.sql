BEGIN;

ALTER TABLE mch_provisioner.servers
DROP CONSTRAINT fk_floating_ip;

ALTER TABLE mch_provisioner.servers
DROP COLUMN addr;

ALTER TABLE mch_provisioner.servers
ADD COLUMN addr INET;

DROP TABLE mch_provisioner.floating_ips;

COMMIT;