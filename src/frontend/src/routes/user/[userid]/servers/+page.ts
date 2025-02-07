export const ssr = false;

import type {PageLoad} from "./$types";
import {env} from "$env/dynamic/public";
import auth from "$lib/service/auth_service";

export const load: PageLoad = async ({ fetch }) => {
  const auth0Client = await auth.createClient();
  const token = await auth0Client.getTokenSilently();

  const res = await fetch(`${env.PUBLIC_BACKEND_URL}/users/1/servers`, {
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  const servers = await res.json();

  let status = 0;
  if (!servers) {
    status = -1;
    return { status };
  }
  let serverHealth = [];
  for (const server of servers) {
    const healthResponse = await fetch(`${env.PUBLIC_BACKEND_URL}/servers/${server.ID}/health`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    const serverip = await healthResponse.json();
    server.ip = serverip.Ip;
    const resStatus = await fetch(`/status/?ip=${server.ip}`)
    const data = await resStatus.json();
    serverHealth.push({data});
  }
  console.log(servers);
  console.log(serverHealth);
  return { servers, serverHealth }; //Packed into an object
};