<!--Check Servers by User-->
<script>
  import {
    Table,
    TableBody,
    TableBodyCell,
    TableBodyRow,
    TableHead,
    TableHeadCell,
  } from "flowbite-svelte";
  let { data } = $props();
</script>

<div class="p-8 mt-16 bg-white dark:bg-gray-900 h-screen">
  <Table>
    <caption
      class="p-5 text-lg font-semibold text-left text-gray-900 bg-white dark:text-white dark:bg-gray-800">
      Your Server(s)
      <p class="mt-1 text-sm font-normal text-gray-500 dark:text-gray-400">
        View and setup your Servers
      </p>
    </caption>
    <TableHead>
      <TableHeadCell>Server Name</TableHeadCell>
      <TableHeadCell>IP</TableHeadCell>
      <TableHeadCell>Version</TableHeadCell>
      <TableHeadCell>Mode</TableHeadCell>
      <TableHeadCell>Difficulty</TableHeadCell>
      <TableHeadCell>Max Players</TableHeadCell>
      <TableHeadCell>PVP</TableHeadCell>
      <TableHeadCell>
        <span class="sr-only">Edit</span>
      </TableHeadCell>
    </TableHead>
    <TableBody tableBodyClass="divide-y">
      {#each data.servers as { id, name, ip, version, mode, difficulty, maxPlayers, pvp }}
        <TableBodyRow>
          <TableBodyCell>{name}</TableBodyCell>
          <TableBodyCell>{ip}</TableBodyCell>
          <TableBodyCell>{version}</TableBodyCell>
          <TableBodyCell>{mode}</TableBodyCell>
          <TableBodyCell>{difficulty}</TableBodyCell>
          <TableBodyCell>{maxPlayers}</TableBodyCell>
          <TableBodyCell>
            {#if pvp}
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
          <TableBodyCell>
            <a
              href="../../server/{id}"
              class="font-medium text-primary-600 hover:underline dark:text-primary-500"
              >Edit</a>
          </TableBodyCell>
        </TableBodyRow>
      {/each}
    </TableBody>
  </Table>
</div>
