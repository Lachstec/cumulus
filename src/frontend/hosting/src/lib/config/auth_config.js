import {
  PUBLIC_AUTH_DOMAIN,
  PUBLIC_AUTH_CLIENT_ID,
  PUBLIC_AUTH_AUDIENCE,
  PUBLIC_AUTH_CACHE_LOCATION,
} from "$env/static/public";

const config = {
  domain: PUBLIC_AUTH_DOMAIN || "ask_auth0",
  clientId: PUBLIC_AUTH_CLIENT_ID || "ask_auth0",
  audience: PUBLIC_AUTH_AUDIENCE || "http://localhost:3001/api", // Optional but needed for access tokens
  cacheLocation: PUBLIC_AUTH_CACHE_LOCATION || "localstorage", // Optional: Use localStorage to persist tokens
};
export default config;
