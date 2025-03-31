export const ssr = false;

import type { PageLoad } from "./$types";
import { env } from "$env/dynamic/public";
import auth from "$lib/service/auth_service";

type UserData = {
  ID: number;
  Sub: "";
  name: string;
  class: string;
};

export const load: PageLoad = async ({ params }) => {
  let data: UserData[];
  let backend_url = env.PUBLIC_BACKEND_URL;

  let auth0Client = await auth.createClient();
  const token = await auth0Client.getTokenSilently();

  const res = await fetch(`${backend_url}/users`, {
    method: "GET",
    headers: {
      Authorization: `Bearer ${token}`,
    },
  });
  data = await res.json();
  return { data };
};
