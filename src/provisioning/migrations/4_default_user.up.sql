BEGIN;

INSERT INTO mch_provisioner.users (sub, name, class) VALUES ('sample_user', 'testuser', 'admin');

COMMIT;