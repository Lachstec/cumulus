import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";
import { PUBLIC_BACKEND_URL} from "$env/static/public";

let backend_url = PUBLIC_BACKEND_URL

export const load: PageLoad = async ({ params }) => {
  const url = `${backend_url}/users/${params.userid}`;

  const res = await fetch(url, {
    method: "GET",
  });

  if (!res.ok) {
    throw error(res.status, "Not found");
  }

  const data = await res.json();

  return {
    name: data.name,
    role: data.role,
  };
};
