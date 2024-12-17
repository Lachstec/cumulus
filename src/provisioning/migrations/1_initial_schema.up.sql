BEGIN;

-- Create a new Schema for the service
CREATE SCHEMA IF NOT EXISTS mch_provisioner;

-- Type that represents the possible states a server can have
CREATE TYPE server_status AS ENUM ('running', 'stopped', 'restarting');

-- Type representing the minecraft game difficulty
CREATE TYPE difficulty AS ENUM ('peaceful', 'easy', 'normal', 'hard');

-- Type representing the minecraft game mode
CREATE TYPE game_mode AS ENUM ('creative', 'adventure', 'survival', 'hardcore');

CREATE TYPE class as ENUM ('admin', 'user');

-- Server table representing gameservers
CREATE TABLE mch_provisioner.servers(
    id SERIAL PRIMARY KEY,         -- id of the server
    openstack_id UUID NOT NULL,    -- UUID in openstack
    userid SERIAL NOT NULL REFERENCES mch_provisioner.users
        ON DELETE CASCADE,
    name VARCHAR(256) NOT NULL,    -- Name of the Server
    addr INET,                     -- IP-Address of the server
    status server_status NOT NULL, -- Current Server Status
    port INTEGER NOT NULL,         -- Port the Server is listening on
    memory_mb INTEGER NOT NULL,    -- Amount of RAM the Server has
    game VARCHAR(128),             -- Which game this server is for
    game_version VARCHAR(128),     -- Which game version is running
    game_mode game_mode,           -- Which game mod is currently active
    difficulty difficulty,         -- Game difficulty
    whitelist_enabled BOOLEAN,     -- Whether the whitelist is enabled
    players_max INTEGER,           -- How many Players are allowed
    ssh_key BYTEA                  -- SSH Key that can be used to connect to the gameserver
);

CREATE TABLE mch_provisioner.users(
    id SERIAL PRIMARY KEY,
    sub VARCHAR(256) NOT NULL,
    name VARCHAR(256) NOT NULL,
    class class NOT NULL,
);

-- Table storing server backup information
CREATE TABLE mch_provisioner.world_backups(
    id SERIAL PRIMARY KEY,                                              -- Id of the backup
    openstack_id UUID NOT NULL,                                         -- UUID in openstack
    server_id SERIAL NOT NULL REFERENCES mch_provisioner.servers
        ON DELETE CASCADE,                                              -- Server the backup belongs to
    timestamp TIMESTAMP,                                                -- Timestamp of creation
    size INTEGER                                                        -- Size of the backup in megabytes
);

COMMIT;