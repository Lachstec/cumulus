<script lang="ts">
  import {
    Indicator,
    Span,
    TableBodyCell,
    TableBodyRow,
  } from "flowbite-svelte";
  import { goto } from "$app/navigation";

  interface Status {
    status: string;
    result: {
      roundTripLatency?: String;
      players?: {
        online?: Number;
      };
    };
  }

  let { server } = $props();
  let stats = $state({ status: "init", result: {} } as Status);
  const updateInterval = 10;

  async function fetchHealth() {
    const resStatus = await fetch(`/status/?ip=${server.ip}`);
    stats = await resStatus.json();
    console.log(`Stats for ${server.ip}:`);
    console.log($state.snapshot(stats));
  }

  $effect(() => {
    const interval = setInterval(fetchHealth, updateInterval * 1000);
    return () => {
      clearInterval(interval);
    };
  });
</script>

<TableBodyRow on:click={() => goto(`../../server/${server.ID}`)}>
  <TableBodyCell class="!p-1">
    <Indicator
      color={stats.status === "success"
        ? "green"
        : stats.status === "error"
          ? "red"
          : "yellow"} />
  </TableBodyCell>
  <TableBodyCell>
    <Span>
      {server.name.length > 20
        ? server.name.substring(0, 20) + "..."
        : server.name}
    </Span>
  </TableBodyCell>
  <TableBodyCell>{server.ip}</TableBodyCell>
  <TableBodyCell
    >{stats.status === "success" ? stats.result.roundTripLatency : "-"} ms</TableBodyCell>
  <TableBodyCell>{server.game_version}</TableBodyCell>
  <TableBodyCell>{server.gamemode}</TableBodyCell>
  <TableBodyCell>{server.difficulty}</TableBodyCell>
  <TableBodyCell
    >{stats.status === "success"
      ? stats.result.players?.online
      : "-"}/{server.players_max}</TableBodyCell>
  <TableBodyCell>
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
        d={server.pvp_enabled
          ? "M5 11.917 9.724 16.5 19 7.5"
          : "M6 18 17.94 6M18 18 6.06 6"} />
    </svg>
  </TableBodyCell>
</TableBodyRow>
