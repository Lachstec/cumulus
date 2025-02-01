<script lang="ts">
  import {
    Toggle,
    Input,
    Label,
    Select,
    Button,
    NumberInput,
    Textarea,
  } from "flowbite-svelte";
  import {
    FloppyDiskSolid,
    TrashBinSolid,
    CaretLeftSolid,
  } from "flowbite-svelte-icons";
  import { PUBLIC_BACKEND_URL } from "$env/static/public";

  // Drop Downs
  let gameDiff = [
    { value: "peaceful", name: "Peaceful" },
    { value: "easy", name: "Easy" },
    { value: "normal", name: "Medium" },
    { value: "hard", name: "Hard" },
  ];

  let gameMode = [
    { value: "survival", name: "Survival" },
    { value: "creative", name: "Creative" },
    { value: "adventure", name: "Adventure" },
    { value: "spectator", name: "Spectator" },
  ];

  // GET
  let { data } = $props();
  let { server, versions } = data;
  let updatedServer = server;
  let whitelistVis = $derived(updatedServer.whitelist_enabled);
  // PATCH
  async function updateServer() {
    console.log(updatedServer);
    console.log("Whitelist:" + updatedServer.whitelist_enabled);
    try {
      const response = await fetch(
        `${PUBLIC_BACKEND_URL}/servers/${server.ID}`,
        {
          method: "PATCH",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify(updatedServer),
        },
      );
      console.log(response);
    } catch (error) {
      console.error(error);
    }
  }

  // DELETE
  async function deleteServer() {
    try {
      const response = await fetch(
        `${PUBLIC_BACKEND_URL}/servers/${server.ID}`,
        {
          method: "DELETE",
          headers: {
            "Content-Type": "application/json",
          },
        },
      );
      const data = await response.json();
      console.log("Update successful:", data);
    } catch (error) {
      console.error(error);
    }
  }
</script>

<form class="p-8 max-w-7xl mx-auto mt-16 bg-white dark:bg-gray-900 h-screen">
  <div class="grid gap-6 mb-6 md:grid-cols-2">
    <div>
      <Label for="server_name" class="mb-2">Server Name</Label>
      <Input
        type="text"
        id="server_name"
        placeholder={updatedServer.name}
        required
        disabled />
    </div>
    <div>
      <Label for="pvp_toggle" class="mb-2">PVP</Label>
      <Toggle bind:checked={updatedServer.pvp_enabled} id="pvp_toggle"></Toggle>
    </div>
    <div>
      <Label>
        Select an difficulty
        <Select
          id="game_difficulty"
          class="mt-2"
          items={gameDiff}
          bind:value={updatedServer.difficulty} />
      </Label>
    </div>
    <div>
      <Label>Select a version</Label>
      <Select
        id="game_version"
        class="mt-2"
        bind:value={updatedServer.game_version}>
        {#each versions as version}
          <option value={version}>{version}</option>
        {/each}
      </Select>
    </div>
    <div>
      <Label>
        Select a Gamemode
        <Select
          id="game_mode"
          class="mt-2"
          items={gameMode}
          bind:value={updatedServer.gamemode} />
      </Label>
    </div>
    <div>
      <Label>
        Maximum Number of Players
        <NumberInput
          id="players_max"
          class="mt-2"
          bind:value={updatedServer.players_max} />
      </Label>
    </div>
    <div>
      <Label for="whitelist_toggle" class="mb-2">Whitelist</Label>
      <Toggle
        bind:checked={updatedServer.whitelist_enabled}
        id="whitelist_toggle"></Toggle>
    </div>
    {#if updatedServer.whitelist_enabled}
      <div>
        <Label for="whitelist_textarea" class="mb-2">Whitelisted Users</Label>
        <Textarea
          id="whitelist_textarea"
          placeholder="Your message"
          rows="4"
          name="whitelist_message" />
      </div>
    {/if}
  </div>
  <Button
    class="mr-1"
    size="md"
    color="yellow"
    href="../user/{server.ID}/servers"
    ><CaretLeftSolid class="w-5 h-5 me-2" />Back</Button>

  <Button class="mr-1" size="md" color="green" on:click={updateServer}
    ><FloppyDiskSolid class="w-5 h-5 me-2" />Save</Button>

  <Button class="mr-1" size="md" color="red" on:click={deleteServer}
    ><TrashBinSolid class="w-5 h-5 me-2" />Delete</Button>
</form>
