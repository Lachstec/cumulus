<script lang="ts">
  import { Toggle, Input, Label, Select, Button } from "flowbite-svelte";
  import { FloppyDiskSolid, TrashBinSolid } from "flowbite-svelte-icons";
  import {PUBLIC_BACKEND_URL, PUBLIC_REQUESTER_NAME} from "$env/static/public";

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
  let {server, versions} = data

  // PATCH
  async function updateServer() {
    console.log(server)
    try {
      const response = await fetch(`${PUBLIC_BACKEND_URL}/servers/${server.ID}`,{
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          name: server.name,
          game: "minecraft",
          game_version: "latest",
          gamemode: "survival",
          difficulty: "normal",
          whitelist_enabled: true,
          pvp_enabled: true,
          players_max: 10,
        }),
      });
    } catch (error) {
      console.error(error);
    }
  }

  // DELETE
  async function deleteServer() {
    try {
      const response = await fetch(`${PUBLIC_BACKEND_URL}/servers/${server.ID}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        }
      });
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
        placeholder="{server.name}"
        required
        disabled
      />
    </div>
    <div>
      <Label for="pvp_toggle" class="mb-2">PVP</Label>
      <Toggle checked={server.pvp_enabled} id="pvp_toggle"></Toggle>
    </div>
    <div>
      <Label>
        Select an difficulty
        <Select
          id="game_difficulty"
          class="mt-2"
          items={gameDiff}
          bind:value={server.difficulty} />
      </Label>
    </div>
    <div>
      <Label>Select a version</Label>
      <Select id="game_version" class="mt-2" bind:value={server.game_version} >
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
          bind:value={server.gamemode} />
      </Label>
    </div>
  </div>
  <Button class="mr-1" size="md" color="green" on:click={updateServer}
    ><FloppyDiskSolid class="w-5 h-5 me-2" />Save</Button>

  <Button class="mr-1" size="md" color="red" on:click={deleteServer}
    ><TrashBinSolid class="w-5 h-5 me-2" />Delete</Button>
</form>
