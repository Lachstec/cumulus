import { servers } from "./data.js";
// läd Serverseitig, weil default alles Client Seitig passiert
// hier als mock
export function load() {
  return {
    servers: servers.map((server) => ({
      id: server.id,
      name: server.name,
      ip: server.ip,
      version: server.version,
      mode: server.mode,
      difficulty: server.difficulty,
      maxPlayers: server.maxPlayers,
      pvp: server.pvp,
    })),
  };
}
