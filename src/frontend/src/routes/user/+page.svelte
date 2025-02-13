<script lang="ts">
  import { onMount } from "svelte";
  import {
    Button,
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
  } from "flowbite-svelte";
  import auth from "$lib/service/auth_service";
  import type { Auth0Client } from "@auth0/auth0-spa-js";
  import { env } from "$env/dynamic/public";
  import type {PageData} from "../../../.svelte-kit/types/src/routes/user/[userid]/$types";

  let backend_url = env.PUBLIC_BACKEND_URL;
  let roleHelper = { "admin":"Admin", "user":"User" };
  type UserData = {
    ID: number;
    Sub: "";
    name: string;
    class: string;
  };

  let {data} = $props();
  let users = data.data as UserData[];
  let isVisible = false;

  function deleteUser(ID: number) {
    fetch(`${backend_url}/users/${ID}`, {
      method: "DELETE",
    })
      .then((res) => {
        if (res.status === 204) {
          isVisible = true;
          //getAllUsers();
        } else {
          console.error("Failed to delete user: ", res.status);
        }
      })
      .catch((error) => {
        console.error(error);
      });
  }
</script>

<ul>
  <div class="p-8 bg-white dark:bg-gray-900">
    <Table>
      <caption
        class="p-5 text-lg font-semibold text-left text-gray-900 bg-white dark:text-white dark:bg-gray-800">
        Users
        <p class="mt-1 text-sm font-normal text-gray-500 dark:text-gray-400">
          View, Edit and Delete Users
        </p>
      </caption>
      <TableHead>
        <TableHeadCell>Username</TableHeadCell>
        <TableHeadCell>Role</TableHeadCell>
        <TableHeadCell>Options</TableHeadCell>
      </TableHead>
      <TableBody tableBodyClass="divide-y">
        {#each users as { name, class: role, ID }}
          <TableBodyRow>
            <TableBodyCell>{name}</TableBodyCell>
            <TableBodyCell>{roleHelper[role as keyof typeof roleHelper]}</TableBodyCell>
            <TableBodyCell>
              <Button size="xs" outline color="green" href="../../user/{ID}"
                >View</Button>
              <Button size="xs" outline color="yellow">Edit</Button>
              <Button
                size="xs"
                outline
                color="red"
                role="button"
                on:click={() => {
                  if (confirm("Do you really want to delete this user?")) {
                    deleteUser(ID);
                    console.log("Deleted user");
                  } else {
                    console.log("Not Deleted");
                  }
                }}>Delete</Button>
            </TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  </div>
</ul>
