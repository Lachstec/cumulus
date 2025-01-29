import type { PageLoad } from './$types';
import {PUBLIC_BACKEND_URL} from "$env/static/public";

export const load: PageLoad = async ({fetch, params}) => {
    // First get Servers from Backend
    const serverRes = await fetch(`${PUBLIC_BACKEND_URL}/servers/${params.serverid}`);
    const server = await serverRes.json();
    //console.log(server);

    // Then get Versions from a 3rd-Party API
    const versionRes = await fetch("https://mc-versions-api.net/api/java");
    const versionData = await versionRes.json();
    const versions = ["latest", ...versionData.result]
    console.log(versions);
    return {
        server,
        versions
    };
}