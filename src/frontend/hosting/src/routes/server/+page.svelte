<script lang="ts">
  import { Card, Button } from "flowbite-svelte";
  import { CheckCircleSolid } from "flowbite-svelte-icons";

  let cards = [
    { ID: 1, title: "Tiny", ram: 512, disk: 1, cpu: 1, cost: 2 },
    { ID: 2, title: "Small", ram: 2048, disk: 20, cpu: 1, cost: 4 },
    { ID: 3, title: "Medium", ram: 4096, disk: 40, cpu: 2, cost: 8 },
    { ID: 4, title: "Large", ram: 8192, disk: 80, cpu: 4, cost: 16 },
    { ID: 5, title: "X-Large", ram: 16384, disk: 160, cpu: 8, cost: 100 },
  ];

  async function orderServer(flavour: number) {
    let response = null;
    try {
      response = await fetch("http://localhost:10000/servers", {
        method: "POST",
        body: JSON.stringify({
          flavour: cards[flavour].ID,
          name: "MyServer",
          image: "d6d1835c-7180-4ca9-b4a1-470afbd8b398",
          game: "minecraft",
          game_version: "latest",
          gamemode: "survival",
          difficulty: "normal",
          whitelist_enabled: true,
          pvp_enabled: true,
          players_max: 20,
        }),
      });
    } catch (err) {
      console.log(err);
    }
    console.log(response);
  }
</script>

<div
  class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 p-4">
  {#each cards as card}
    <Card padding="xl">
      <h5 class="mb-4 text-xl font-medium text-gray-500 dark:text-gray-400">
        {card.title}
      </h5>
      <div class="flex items-baseline text-gray-900 dark:text-white">
        <span class="text-5xl font-extrabold tracking-tight">{card.cost}</span>
        <span class="text-3xl font-semibold">€</span>
        <span class="ms-1 text-xl font-normal text-gray-500 dark:text-gray-400"
          >/month</span>
      </div>
      <!-- List -->
      <ul class="my-7 space-y-4">
        <li class="flex space-x-2 rtl:space-x-reverse">
          <CheckCircleSolid
            class="w-4 h-4 text-primary-600 dark:text-primary-500" />
          <span
            class="text-base font-normal leading-tight text-gray-500 dark:text-gray-400">
            {card.ram} MB RAM
          </span>
        </li>
        <li class="flex space-x-2 rtl:space-x-reverse">
          <CheckCircleSolid
            class="w-4 h-4 text-primary-600 dark:text-primary-500" />
          <span
            class="text-base font-normal leading-tight text-gray-500 dark:text-gray-400">
            {card.disk} GB Data
          </span>
        </li>
        <li class="flex space-x-2 rtl:space-x-reverse">
          <CheckCircleSolid
            class="w-4 h-4 text-primary-600 dark:text-primary-500" />
          <span
            class="text-base font-normal leading-tight text-gray-500 dark:text-gray-400">
            {card.cpu} vCPUs</span>
        </li>
      </ul>
      <Button class="w-full" on:click={() => orderServer(card.ID - 1)}
        >Choose Flavour</Button>
    </Card>
  {/each}
</div>
