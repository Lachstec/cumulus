<script>
  // @ts-nocheck
    import { onMount } from "svelte";
    import { get } from 'svelte/store';
    // @ts-ignore
    import auth from "../auth_service";
    import { isAuthenticated, user } from "$lib/store/stores";
  
    // @ts-ignore
    let auth0Client;
    // @ts-ignore
    let newTask;
  
    onMount(async () => {
      console.log("onMountCalled")
      
      auth0Client = await auth.createClient();
  
      isAuthenticated.set(await auth0Client.isAuthenticated());
      user.set(await auth0Client.getUser());
      console.log(get(isAuthenticated))
      console.log(get(user))
    });
  
    function login() {
      auth.loginPopup(auth0Client);
    }
  
    function logout() {
      auth.logout(auth0Client);
    }
</script>

<style>
    #main-application {
      margin-top: 50px;
    }
</style>


<div class="mt-16">
    <!-- App Bar -->
    <nav class="navbar navbar-expand-lg navbar-dark bg-dark">
      <a class="navbar-brand" href="/#">Task Manager</a>
      <button
        class="navbar-toggler"
        type="button"
        data-toggle="collapse"
        data-target="#navbarText"
        aria-controls="navbarText"j
        aria-expanded="false"
        aria-label="Toggle navigation"
      >
        <span class="navbar-toggler-icon"></span>
      </button>
        <span class="navbar-text">
          <ul class="navbar-nav float-right">
            {#if !get(isAuthenticated)}
            <a class="nav-link" href="/#" on:click="{login}">"Log In" </a>
            {:else}
            <a class="nav-link" href="/#" on:click="{logout}">"Log Out" </a>
            {/if}
          </ul>
        </span>
    </nav>
  </div>