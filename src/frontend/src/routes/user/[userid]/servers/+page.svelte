<!--Check Servers by User-->
<script>
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
    Indicator,
    Span,
    P,
    Progressbar,
  } from "flowbite-svelte";
  import { sineInOut } from "svelte/easing";
  import { goto } from "$app/navigation";
  let { data } = $props();
  let indicatorColor = "black";
  let latestData = $state([...data.serverHealth]); // Make a reactive copy
  let progress = $state(0);
  const updateInterval = 10; // interval in secondes

  async function fetchHealth() {
    let serverHealth = [];

    for (let i = 0; i < data.servers.length; i++) {
      const server = data.servers[i];
      const resStatus = await fetch(`/status/?ip=${server.ip}`);
      const health = await resStatus.json();
      serverHealth.push(health);
    }

    // Update each item reactively
    latestData.forEach((_, i) => {
      latestData[i] = serverHealth[i];
    });

    //console.log("Updated health:", latestData);
  }

  $effect(() => {
    const interval = setInterval(fetchHealth, updateInterval * 1000);
    const progressInterval = setInterval(() => {
      progress = (progress + 100 / updateInterval) % 101;
    }, 1000);
    return () => {
      clearInterval(interval);
      clearInterval(progressInterval);
    };
  });
</script>

<div class="p-8 bg-white dark:bg-gray-900">
  {#if data.status === -1}
    <P>You currently do not have any servers</P>
  {:else}
    <Table hoverable="true">
      <caption
        class="p-5 text-lg font-semibold text-left text-gray-900 bg-white dark:text-white dark:bg-gray-800">
        Your Server(s)
        <div class="flex justify-between items-center w-full">
          <span class="text-sm font-normal text-gray-500 dark:text-gray-400">
            View and setup your Servers. To edit a Server, just click on it.
          </span>
          <Progressbar
            class="w-40"
            {progress}
            size="h-4"
            animate
            tweenDuration={500}
            easing={sineInOut} />
        </div>
      </caption>
      <TableHead>
        <TableHeadCell class="!p-1"></TableHeadCell>
        <TableHeadCell>Server Name</TableHeadCell>
        <TableHeadCell>IP</TableHeadCell>
        <TableHeadCell>Ping</TableHeadCell>
        <TableHeadCell>Version</TableHeadCell>
        <TableHeadCell>Mode</TableHeadCell>
        <TableHeadCell>Difficulty</TableHeadCell>
        <TableHeadCell>Players</TableHeadCell>
        <TableHeadCell>PVP</TableHeadCell>
      </TableHead>
      <TableBody tableBodyClass="divide-y">
        {#each data.servers as { ID, Status, name, ip, game_version, gamemode, difficulty, players_max, pvp_enabled }, index}
          <TableBodyRow on:click={() => goto(`../../server/${ID}`)}>
            <TableBodyCell class="!p-1">
              <Indicator color={Status === "running" ? "green" : "red"} />
            </TableBodyCell>
            <TableBodyCell>
              <Span>
                {name.length > 20 ? name.substring(0, 20) + "..." : name}
              </Span>
            </TableBodyCell>
            <TableBodyCell>{ip}</TableBodyCell>
            <TableBodyCell
              >{latestData[index]?.roundTripLatency}ms</TableBodyCell>
            <TableBodyCell>{game_version}</TableBodyCell>
            <TableBodyCell>{gamemode}</TableBodyCell>
            <TableBodyCell>{difficulty}</TableBodyCell>
            <TableBodyCell
              >{latestData[index]?.players.online}/{players_max}</TableBodyCell>
            <TableBodyCell>
              {#if pvp_enabled}
                <svg
                  class="w-6 h-6 text-gray-800 dark:text-white"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  fill="none"
                  viewBox="0 0 24 24">
                  <path
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M5 11.917 9.724 16.5 19 7.5" />
                </svg>
              {:else}
                <svg
                  class="w-6 h-6 text-gray-800 dark:text-white"
                  aria-hidden="true"
                  xmlns="http://www.w3.org/2000/svg"
                  width="24"
                  height="24"
                  fill="none"
                  viewBox="0 0 24 24">
                  <path
                    stroke="currentColor"
                    stroke-linecap="round"
                    stroke-linejoin="round"
                    stroke-width="2"
                    d="M6 18 17.94 6M18 18 6.06 6" />
                </svg>
              {/if}
            </TableBodyCell>
          </TableBodyRow>
        {/each}
      </TableBody>
    </Table>
  {/if}
</div>
