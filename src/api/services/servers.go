package services

import (
	Data "data"
)

func ReadNumOfServers() int {
	return len(Data.Servers)
}

func ReadAllServers() []Data.Server {
	return Data.Servers
}

func ReadServerByServerID(serverid int) Data.Server {
	return Data.Servers[serverid]
}

func CreateServer(server Data.Server) {
	Data.Servers = append(Data.Servers, server)
}

func DeleteServerByServerID(serverid int) {
	Data.Servers = append(Data.Servers[:serverid], Data.Servers[serverid + 1:]...)
}