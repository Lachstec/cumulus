BEGIN;

DELETE FROM mch_provisioner.users WHERE name = 'testuser';

COMMIT;