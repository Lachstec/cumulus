BEGIN;

CREATE TABLE mch_provisioner.floating_ips(
    id SERIAL PRIMARY KEY,
    openstack_id UUID NOT NULL,
    addr INET NOT NULL
);

ALTER TABLE mch_provisioner.servers
DROP COLUMN addr;

ALTER TABLE mch_provisioner.servers
ADD COLUMN addr INT;

ALTER TABLE mch_provisioner.servers
ADD CONSTRAINT fk_floating_ip
FOREIGN KEY (addr)
REFERENCES mch_provisioner.floating_ips(id);

COMMIT;