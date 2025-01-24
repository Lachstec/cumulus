import {
  PUBLIC_AUTH_DOMAIN,
  PUBLIC_AUTH_CLIENT_ID,
  PUBLIC_AUTH_AUDIENCE,
  PUBLIC_AUTH_CACHE_LOCATION,
} from "$env/static/public";

const config = {
  domain: PUBLIC_AUTH_DOMAIN,
  clientId: PUBLIC_AUTH_CLIENT_ID,
  audience: PUBLIC_AUTH_AUDIENCE, // Optional but needed for access tokens
  cacheLocation: PUBLIC_AUTH_CACHE_LOCATION, // Optional: Use localStorage to persist tokens
};
export default config;
