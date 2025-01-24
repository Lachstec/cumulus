<script lang="ts">
  import {
    Button,
    Modal,
    Label,
    Input,
    Select,
    AccordionItem,
    Accordion,
  } from "flowbite-svelte";
  import { ExclamationCircleOutline } from "flowbite-svelte-icons";
  import type { PageData } from "./$types";
  import { page } from "$app/stores";
  import { PUBLIC_BACKEND_URL } from "$env/static/public";

  let backend_url = PUBLIC_BACKEND_URL;

  let { data }: { data: PageData } = $props();

  let showModal = $state(false);
  let isButtonEnabled = $state(false);
  let timer = $state(5);

  function startTimer() {
    isButtonEnabled = false;
    timer = 5;

    const countdown = setInterval(() => {
      timer--;
      if (timer <= 0) {
        clearInterval(countdown);
        isButtonEnabled = true;
      }
    }, 1000);
  }

  function openModal() {
    showModal = true;
    startTimer();
  }

  function closeModal() {
    showModal = false;
  }
  let selectedName = $state(data.name);
  let selectedRole = $state(data.role);
  let roles = [
    { value: "Admin", name: "admin" },
    { value: "User", name: "user" },
  ];

  async function updateData() {
    console.log("Updated Data: " + selectedName + " " + selectedRole);
    try {
      const userId = $page.url.pathname.split("/").pop(); // Assuming the ID is in the URL
      if (!userId) {
        throw new Error("User ID not found in the URL");
      }

      const response = await fetch(
        `${backend_url}/users/${$page.params.userid}`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            name: selectedName,
            role: selectedRole,
          }),
        },
      );

      if (response.ok) {
        console.log("User updated successfully");
      } else {
        throw new Error(`Failed to update user: ${response.statusText}`);
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function deleteUser() {
    try {
      const response = await fetch(
        `${backend_url}/users/${$page.params.userid}`,
        { method: "DELETE" },
      );
      if (response.ok) {
        console.log("Item deleted successfully");
      } else {
        throw new Error("Failed to delete");
      }
    } catch (err) {
      console.error(err);
    }
  }
</script>

<div class="p-8 mt-16 mb-6 bg-white dark:bg-gray-900 h-screen">
  <div class="mb-6">
    <Label for="name-input" class="block mb-2">Name</Label>
    <Input id="name-input" bind:value={selectedName} />
  </div>

  <div class="mb-6">
    <Label>
      Select a role
      <Select class="mt-2" items={roles} bind:value={selectedRole} />
    </Label>
  </div>

  <Button class="mb-6" color="green" on:click={updateData}>Save</Button>

  <!-- Danger Zone -->
  <Accordion>
    <AccordionItem>
      <span slot="header">Danger Zone</span>
      <Button on:click={openModal} color="red">Delete Account</Button>
    </AccordionItem>
  </Accordion>
</div>

<Modal bind:open={showModal} on:close={closeModal} size="xs" autoclose>
  <div class="text-center">
    <ExclamationCircleOutline
      class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
    <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
      Are you sure you want to delete your account?
    </h3>
    {#if timer > 0}
      <p class="mb-5">Please wait for: {timer}s</p>
    {/if}
    <Button
      color="red"
      class="me-2"
      disabled={!isButtonEnabled}
      on:click={deleteUser}>Yes, I'm sure</Button>
    <Button color="alternative">No, cancel</Button>
  </div>
</Modal>
