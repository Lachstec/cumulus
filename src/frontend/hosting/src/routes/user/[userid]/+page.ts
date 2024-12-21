import { error } from "@sveltejs/kit";
import type { PageLoad } from "./$types";

export const load: PageLoad = async ({ params }) => {
  const url = `http://localhost:10000/users/${params.userid}`;

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
