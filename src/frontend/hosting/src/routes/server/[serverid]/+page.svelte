<script lang="ts">
  import { Toggle, Input, Label, Helper, Select } from "flowbite-svelte";
  import { onMount } from "svelte";
  let selected_gamemode: String;
  let selected_version: String;
  let gamemode = [
    { value: "peaceful", name: "Peaceful" },
    { value: "easy", name: "Easy" },
    { value: "medium", name: "Medium" },
    { value: "hard", name: "Hard" },
  ];

  let versions: String[] = ["latest"];
  onMount(async () => {
    const res = await fetch("https://mc-versions-api.net/api/java");
    const data = await res.json();
    versions = versions.concat(data.result);
    console.log(versions);
  });
</script>

<form class="p-8 max-w-7xl mx-auto mt-16 bg-white dark:bg-gray-900 h-screen">
  <div class="grid gap-6 mb-6 md:grid-cols-2">
    <div>
      <Label for="server_name" class="mb-2">Server Name</Label>
      <Input type="text" id="server_name" placeholder="Bort" required />
    </div>
    <div>
      <Label for="pvp_toggle" class="mb-2">PVP</Label>
      <Toggle checked={true} id="pvp_togglge"></Toggle>
    </div>
    <div>
      <Label>
        Select an difficulty
        <Select
          id="game_difficulty"
          class="mt-2"
          items={gamemode}
          bind:value={selected_gamemode} />
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
  </div>
</form>
