import { env } from "$env/dynamic/public";

const config = {
  domain: env.PUBLIC_AUTH_DOMAIN,
  clientId: env.PUBLIC_AUTH_CLIENT_ID,
  audience: env.PUBLIC_AUTH_AUDIENCE, // Optional but needed for access tokens
  cacheLocation: env.PUBLIC_AUTH_CACHE_LOCATION, // Optional: Use localStorage to persist tokens
};
export default config;
