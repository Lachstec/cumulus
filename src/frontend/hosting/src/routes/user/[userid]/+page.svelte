<script lang="ts">
    import { Button, Modal, Label, Input, Select, AccordionItem, Accordion } from "flowbite-svelte";
    import { ExclamationCircleOutline } from "flowbite-svelte-icons";
    import type { PageData } from './$types';

    let { data }: { data: PageData } = $props();

    let showModal = $state(false);
    let isButtonEnabled = $state(false);
    let timer = $state(5);

    function startTimer() {
        isButtonEnabled = false;
        timer = 5;

        const countdown = setInterval(() => {
            timer--;
            if (timer <= 0) {
                clearInterval(countdown);
                isButtonEnabled = true;
            }
        }, 1000);
    }

    function openModal() {
        showModal = true;
        startTimer();
    }

    function closeModal() {
        showModal = false;
    }

    let selected;
    let roles = [
        { value: "Admin", name: "Admin" },
        { value: "User", name: "User" },
    ]
</script>

<div class="p-8 mt-16 mb-6 bg-white dark:bg-gray-900 h-screen">
    <form method="POST">
        <div class="mb-6">
            <Label for="name-input" class="block mb-2">Name</Label>
            <Input id="name-input" value="{data.name}" />
        </div>
        <div class="mb-6">
        <Label>
            Select a role
            <Select class="mt-2" items={roles} bind:value={data.role} />
        </Label>
        </div>

        <Button class="mb-6" color="green">Save</Button>
    </form>


    <Accordion>
        <AccordionItem>
            <span slot="header">Danger Zone</span>
            <Button on:click={openModal} color="red">Delete Account</Button>
        </AccordionItem>
    </Accordion>
</div>

<Modal bind:open={showModal} on:close={closeModal} size="xs" autoclose>
    <div class="text-center">
        <ExclamationCircleOutline
                class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
        <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">
            Are you sure you want to delete your account?
        </h3>
        {#if timer > 0}
            <p class="mb-5">Please wait for: {timer}s</p>
        {/if}
        <Button color="red" class="me-2" disabled={!isButtonEnabled}
        >Yes, I'm sure</Button>
        <Button color="alternative">No, cancel</Button>
    </div>
</Modal>