<!--Check Servers by User-->
<script>
  import {
    Table,
    TableBody,
    TableHead,
    TableHeadCell,
    P,
    Progressbar,
  } from "flowbite-svelte";
  import { sineInOut } from "svelte/easing";
  import Row from "./Row.svelte";
  let { data } = $props();
  let progress = $state(0);
  const updateInterval = 10; // interval in secondes

  $effect(() => {
    const progressInterval = setInterval(() => {
      progress = (progress + 100 / updateInterval) % 101;
    }, 1000);
    return () => {
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
            <br />
            This page is continuously updating the server status.
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
        {#each data.servers as server}
          <Row {server} />
        {/each}
      </TableBody>
    </Table>
  {/if}
</div>
