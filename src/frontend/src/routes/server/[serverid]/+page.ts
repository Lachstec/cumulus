import type { PageLoad } from "./$types";
import { env } from "$env/dynamic/public";

export const load: PageLoad = async ({ fetch, params }) => {
  // First get Servers from Backend
  const serverRes = await fetch(
    `${env.PUBLIC_BACKEND_URL}/servers/${params.serverid}`,
  );
  const server = await serverRes.json();
  //console.log(server);

  // Then get Versions from a 3rd-Party API
  const versionRes = await fetch("https://mc-versions-api.net/api/java");
  const versionData = await versionRes.json();
  const versions = ["latest", ...versionData.result];
  console.log(versions);
  return {
    server,
    versions,
  };
};
