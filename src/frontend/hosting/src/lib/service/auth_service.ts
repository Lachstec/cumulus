import {
  Auth0Client,
  createAuth0Client,
  type Auth0ClientOptions,
  type PopupLoginOptions,
} from "@auth0/auth0-spa-js";
import { user, isAuthenticated, popupOpen } from "$lib/store/auth_store";
import config from "$lib/config/auth_config";

async function createClient() {
  let auth0Client = await createAuth0Client({
    domain: config.domain,
    clientId: config.clientId,
  });

  return auth0Client;
}

async function loginPopup(client: Auth0Client, options?: PopupLoginOptions) {
  popupOpen.set(true);
  try {
    await client.loginWithPopup(options);
    user.set(await client.getUser());
    isAuthenticated.set(true);
  } catch (e) {
    console.error(e);
  } finally {
    popupOpen.set(false);
  }
}

function logout(client: Auth0Client) {
  client.logout();
}

const auth = {
  createClient,
  loginPopup,
  logout,
};

export default auth;
