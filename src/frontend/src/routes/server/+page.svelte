<script lang="ts">
  import { Alert, Button, Card, Modal, Spinner } from "flowbite-svelte";
  import { CheckCircleSolid } from "flowbite-svelte-icons";
  import { v4 as uuidv4 } from "uuid";
  import { env } from "$env/dynamic/public";
  import auth from "$lib/service/auth_service";

  let backend_url = env.PUBLIC_BACKEND_URL;

  let cards = [
    { ID: 1, title: "Tiny", ram: 512, disk: 1, cpu: 1, cost: 2 },
    { ID: 2, title: "Small", ram: 2048, disk: 20, cpu: 1, cost: 4 },
    { ID: 3, title: "Medium", ram: 4096, disk: 40, cpu: 2, cost: 8 },
    { ID: 4, title: "Large", ram: 8192, disk: 80, cpu: 4, cost: 16 },
    { ID: 5, title: "X-Large", ram: 16384, disk: 160, cpu: 8, cost: 100 },
  ];

  // While waiting for Response
  let isLoading = false;
  let modalOpen = false;

  //alerts
  let responseOK = false;
  let responseError = false;
  let errorMsg = "generic Error";

  async function orderServer(flavour: number) {
    //Auth
    const auth0Client = await auth.createClient();
    const token = await auth0Client.getTokenSilently();

    isLoading = true;
    modalOpen = true;
    let response = null;
    // Name needs to be unique
    const uuid = uuidv4();
    // Ability to track who ordered a server
    try {
      response = await fetch(`${backend_url}/servers`, {
        method: "POST",
        headers: {
          Authorization: `Bearer ${token}`,
        },
        body: JSON.stringify({
          flavour: cards[flavour].ID,
          name: env.PUBLIC_REQUESTER_NAME + "_" + uuid,
          image: "d6d1835c-7180-4ca9-b4a1-470afbd8b398",
          game: "minecraft",
          game_version: "latest",
          gamemode: "survival",
          difficulty: "normal",
          whitelist_enabled: true,
          pvp_enabled: false,
          players_max: 20,
        }),
      });
    } catch (err) {
      console.log(err);
    }
    console.log(response);
    isLoading = false;
    modalOpen = false;

    if (response?.status === 200) {
      responseOK = true;
    } else {
      responseError = true;
      errorMsg = await response?.text();
    }
  }
</script>

<div class="p-8 bg-white dark:bg-gray-900">
  {#if responseOK}
    <Alert color="green" on:close={() => (responseOK = false)}>
      <span class="font-medium">Success!</span> Your server has been created.
    </Alert>
  {/if}
  {#if responseError}
    <Alert color="red" on:close={() => (responseError = false)}>
      <span class="font-medium">Error!</span> Something went wrong. Error: {errorMsg}
    </Alert>
  {/if}
  <div
    class="grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 xl:grid-cols-5 gap-6 p-4">
    {#each cards as card}
      <Card padding="xl">
        <h5 class="mb-4 text-xl font-medium text-gray-500 dark:text-gray-400">
          {card.title}
        </h5>
        <div class="flex items-baseline text-gray-900 dark:text-white">
          <span class="text-5xl font-extrabold tracking-tight"
            >{card.cost}</span>
          <span class="text-3xl font-semibold">€</span>
          <span
            class="ms-1 text-xl font-normal text-gray-500 dark:text-gray-400"
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
        <Button
          class="w-full"
          on:click={() => orderServer(card.ID - 1)}
          disabled={isLoading}>Choose Flavour</Button>
      </Card>
    {/each}
  </div>
</div>
<Modal
  title="Please Wait, Creating your Server"
  bind:open={modalOpen}
  size="xs">
  <div class="flex flex-col items-center justify-center text-center space-y-4">
    <Spinner />
    <p class="text-base leading-relaxed text-gray-500 dark:text-gray-400">
      This window will close when creation is finished.
    </p>
  </div>
</Modal>
