<script>
    import {Indicator, Span, TableBodyCell, TableBodyRow} from "flowbite-svelte";
    import {goto} from "$app/navigation";
    let { server, health } = $props()
</script>


<TableBodyRow on:click={() => goto(`../../server/${server.ID}`)}>
    <TableBodyCell class="!p-1">
        <Indicator color={health.status==="success" ? "green" : health.status==="error" ? "red" : "yellow"}/>
    </TableBodyCell>
    <TableBodyCell>
          <Span>
            {server.name.length > 20 ? server.name.substring(0, 20) + "..." : server.name}
          </Span>
    </TableBodyCell>
    <TableBodyCell>{server.ip}</TableBodyCell>
    <TableBodyCell
    >{health.status==="success" ? health.result.roundTripLatency : "-" } ms</TableBodyCell>
    <TableBodyCell>{server.game_version}</TableBodyCell>
    <TableBodyCell>{server.gamemode}</TableBodyCell>
    <TableBodyCell>{server.difficulty}</TableBodyCell>
    <TableBodyCell
    >{health.status==="success" ? health.result.players.online : "-" }/{server.players_max}</TableBodyCell>
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
                    d={server.pvp_enabled ?"M5 11.917 9.724 16.5 19 7.5" : "M6 18 17.94 6M18 18 6.06 6"} />
        </svg>
    </TableBodyCell>
</TableBodyRow>
