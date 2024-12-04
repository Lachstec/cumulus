<script lang="ts">
  import {
    Avatar,
    DropdownDivider,
    DropdownHeader,
    DropdownItem,
    NavHamburger,
    Dropdown,
  } from "flowbite-svelte";

  // @ts-nocheck
  import { onMount } from "svelte";
  import auth from "$lib/service/auth_service";
  import { isAuthenticated, user } from "$lib/store/auth_store";
  import type { Auth0Client } from "@auth0/auth0-spa-js";

  let auth0Client: Auth0Client;

  onMount(async () => {
    console.log("onMountCalled");

    auth0Client = await auth.createClient();

    isAuthenticated.set(await auth0Client.checkSession());
    user.set(await auth0Client.getUser());

  });

  async function login() {
    await auth.loginPopup(auth0Client);
    isAuthenticated.set(await auth0Client.isAuthenticated()); // Update authentication status
    user.set(await auth0Client.getUser()); // Update user data
  }

  async function logout() {
    await auth.logout(auth0Client);
    isAuthenticated.set(false); // Reset authentication status
    user.set(null); // Reset user data
  }

  async function getToken() {
    try {
      const token = await auth0Client.getTokenSilently();
      console.log('Token:', token);
    } catch (error) {
      console.error('Error getting token:', error);
    }
  }

  const authStatus = $isAuthenticated;
  const userData = $user;
</script>
<button on:click={getToken}>Get Token</button>
{#if $isAuthenticated}
  <div class="flex items-center space-x-4" id="user-menu">
    <Avatar src={$user?.picture} />
    <div class="space-y-1 font-medium dark:text-white">{$user?.nickname}</div>
    <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
  </div>
  <Dropdown placement="bottom" triggeredBy="#user-menu">
    <DropdownHeader>
      <span class="block text-sm">{$user?.name}</span>
      <span class="block truncate text-sm font-medium">{$user?.email}</span>
    </DropdownHeader>
    <DropdownItem>Dashboard</DropdownItem>
    <DropdownItem href="/user/1">Settings</DropdownItem>
    <DropdownItem href="/user/1/servers">My Servers</DropdownItem>
    <DropdownDivider />
    <DropdownItem
      ><a class="nav-link" href="/#" on:click={logout}>Sign Out</a
      ></DropdownItem>
  </Dropdown>
{:else}
  <div class="flex items-center space-x-4">
    <Avatar />
    <div class="space-y-1 font-medium dark:text-white">
      <a class="nav-link" href="/#" on:click={login}>Sign In</a>
    </div>
    <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
  </div>
{/if}
