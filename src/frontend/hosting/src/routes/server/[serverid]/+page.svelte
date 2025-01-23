<script lang="ts">
  import { Toggle, Input, Label, Select, Button } from "flowbite-svelte";
  import { FloppyDiskSolid } from "flowbite-svelte-icons";
  import { onMount } from "svelte";
  let selected_gameMode: String;
  let selected_gameDiff: String;
  let selected_version: String;
  let gameDiff = [
    { value: "peaceful", name: "Peaceful" },
    { value: "easy", name: "Easy" },
    { value: "medium", name: "Medium" },
    { value: "hard", name: "Hard" },
  ];

  let gameMode = [
    { value: "survival", name: "Survival" },
    { value: "creative", name: "Creative" },
    { value: "adventure", name: "Adventure" },
    { value: "spectator", name: "Spectator" },
  ];

  let versions: String[] = ["latest"];
  onMount(async () => {
    const res = await fetch("https://mc-versions-api.net/api/java");
    const data = await res.json();
    versions = versions.concat(data.result);
  });

  async function sendSettings() {
    console.log("Button")
  }
</script>

<form class="p-8 max-w-7xl mx-auto mt-16 bg-white dark:bg-gray-900 h-screen">
  <div class="grid gap-6 mb-6 md:grid-cols-2">
    <div>
      <Label for="server_name" class="mb-2">Server Name</Label>
      <Input type="text" id="server_name" placeholder="Placeholder Name" required />
    </div>
    <div>
      <Label for="pvp_toggle" class="mb-2">PVP</Label>
      <Toggle checked={true} id="pvp_toggle"></Toggle>
    </div>
    <div>
      <Label>
        Select an difficulty
        <Select
          id="game_difficulty"
          class="mt-2"
          items={gameDiff}
          bind:value={selected_gameDiff} />
      </Label>
    </div>
    <div>
      <Label>Select a version</Label>
      <Select id="game_version" class="mt-2" bind:value={selected_version}>
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
                bind:value={selected_gameMode} />
      </Label>
    </div>
  </div>
  <Button size="md" color="green" on:click = { () => sendSettings }><FloppyDiskSolid class="w-5 h-5 me-2" />Save</Button>
</form>
