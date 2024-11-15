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
  import { get } from "svelte/store";
  import auth from "$lib/service/auth_service";
  import { isAuthenticated, user } from "$lib/store/auth_store";
  import type { Auth0Client } from "@auth0/auth0-spa-js";

  let auth0Client: Auth0Client;
  // @ts-ignore
  let newTask;

  onMount(async () => {
    console.log("onMountCalled");

    auth0Client = await auth.createClient();

    isAuthenticated.set(await auth0Client.isAuthenticated());
    user.set(await auth0Client.getUser());
    console.log(get(isAuthenticated));
    const userData = await auth0Client.getUser();
    console.log(userData);
  });

  function login() {
    auth.loginPopup(auth0Client);
  }

  function logout() {
    auth.logout(auth0Client);
  }
</script>

{#if get(isAuthenticated)}
  <div class="flex items-center space-x-4" id="user-menu">
    <Avatar src="/images/example-profile.png" />
    <div class="space-y-1 font-medium dark:text-white">Username</div>
    <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
  </div>
  <Dropdown placement="bottom" triggeredBy="#user-menu">
    <DropdownHeader>
      <span class="block text-sm">Beep Boop</span>
      <span class="block truncate text-sm font-medium">name@example.com</span>
    </DropdownHeader>
    <DropdownItem>Dashboard</DropdownItem>
    <DropdownItem>Settings</DropdownItem>
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
