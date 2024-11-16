<script>
    import { Button, Modal } from 'flowbite-svelte';
    import { ExclamationCircleOutline } from 'flowbite-svelte-icons';
    let popupModal = false;

    let buttonDeleteDisabled = true;
    let counter = 5;

    function trigger() {
        buttonDeleteDisabled = true;
        counter = 5;
        timeout()
    }
    function timeout(){
        if (--counter > 0) return setTimeout(timeout, 1000);
        buttonDeleteDisabled = false;
        console.log(buttonDeleteDisabled)
    }
    $: if (popupModal) {
        trigger();
    }
</script>

<Button on:click={() => (popupModal = true)} color="red">Delete Account</Button>

<Modal bind:open={popupModal} size="xs" autoclose>
    <div class="text-center">
        <ExclamationCircleOutline class="mx-auto mb-4 text-gray-400 w-12 h-12 dark:text-gray-200" />
        <h3 class="mb-5 text-lg font-normal text-gray-500 dark:text-gray-400">Are you sure you want to delete your account?</h3>
        {#if counter > 0}
            <p class="mb-5">Please wait for: {counter}s</p>
        {/if}
        <Button color="red" class="me-2" disabled={buttonDeleteDisabled}>Yes, I'm sure</Button>
        <Button color="alternative">No, cancel</Button>
    </div>
</Modal>