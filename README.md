# Cumulus - Gameservers in the Cloud
This is the repository for cumulus, a web platform to deploy gameserves
in the cloud by leveraging the OpenStack cloud platform. Cumulus serves as 
a graded project for the Course "Cloud Services" at Fulda University of Applied Sciences.

### Participants
- János Euler --> Frontend, Auth
- Moritz Freund --> API, Infrastructure
- Leon Lux --> Backend, Database, DevOps

## Getting Started
Cumulus consists of three components:
- Frontend Webapplication
- Backend API
- SQL Database

For local testing and development, you can simply clone the repository and 
work on the components separately. The backend consists of a single Go Application. You need a 
local PostgreSQL Database in order to run it. You can spawn one with docker by running this 
command: 

```bash
docker run -p5432:5432 -e POSTGRES_HOST_AUTH_METHOD=trust  -d postgres:latest
```
The backend can be spawned by runnign this command in `src/provisioning`:

```bash
go mod tidy
go run cmd/provisioner/main.go
```

You can then reach the backend on `http://localhost:10000`. An overview of the 
existing endpoints gets printed to the command line.

TODO: Frontend description

## Architecture Overview
TODO

## Deployment
To deploy cumulus, you can use the provided Terraform files. It creates three compute instances, one for the backend,
one for the frontend and one for the database. This offers a simple and effective way to get the project running,
but it does not provide redundancy for the database. It is therefore recommended to periodically dump the 
database and store it in a safe and fault-tolerant storage, like OpenStack Object Storage.

The provided Terraform setup needs an OpenStack cluster as deployment target. The
Cluster also needs to provide the following Services:
- Nova (Compute)
- Cinder (Block Storage)
- Neutron (Networking)
- Optional: Designate (DNS)

Furthermore, cumulus requires an account that has API access to the
aforementioned services.

## Montioring / Logging
TODO

## Troubleshooting
TODO

## Service / Maintenance
TODO

## Auth, Accounting, Permissionmanagement
TODO

## User Groups
TODO?

## Backup / Scaling / Disaster Recovery / Archiving
TODO
