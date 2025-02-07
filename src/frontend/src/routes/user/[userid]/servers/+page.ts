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
  const ip = '10.32.6.17'
  const response = await fetch(`/status/?ip=${ip}`);
  console.log(await response.json());

  for (const server of servers) {
    const healthResponse = await fetch(`${env.PUBLIC_BACKEND_URL}/servers/${server.ID}/health`, {
      headers: {
        Authorization: `Bearer ${token}`,
      },
    })
    const serverip = await healthResponse.json();
    console.log(serverip[0]);
    server.ip = serverip[0].Ip;
  }
  console.log(servers);
  return { servers }; //Packed into an object
};