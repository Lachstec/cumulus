import type { PageLoad } from "./$types";
import { PUBLIC_BACKEND_URL } from "$env/static/public";

export const load: PageLoad = async ({ fetch, params }) => {
  console.log(PUBLIC_BACKEND_URL);
  const res = await fetch(`${PUBLIC_BACKEND_URL}/users/1/servers`);
  const servers = await res.json();
  return { servers }; //Packed into an object
};
