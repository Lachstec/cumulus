import type { PageLoad } from "./$types";
import { PUBLIC_BACKEND_URL } from "$env/static/public";
import auth from "$lib/service/auth_service";


export const load: PageLoad = async ({ fetch }) => {
  const auth0Client = await auth.createClient();
  const token = await auth0Client.getTokenSilently();

  const res = await fetch(`${PUBLIC_BACKEND_URL}/users/1/servers`,{
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  const servers = await res.json();
  return { servers }; //Packed into an object
};
