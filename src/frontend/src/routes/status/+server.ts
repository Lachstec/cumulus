import { json } from "@sveltejs/kit";
import type { RequestHandler } from './$types';

export const GET: RequestHandler = async ({url}) =>{
    const util = await import('minecraft-server-util'); // ✅ Dynamic import for CommonJS
    const ip = String(url.searchParams.get("ip"));
    const options = {
        timeout: 1000 * 5, // 5 seconds
        enableSRV: true // Enable SRV record lookup
    };

    try {
        const result = await util.status(ip, 25565, options);
        return json(result);
    } catch (error) {
        console.error("Error fetching server status:", error);
        return json(error);
    }
}