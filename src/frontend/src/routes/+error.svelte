<script lang="ts">
    import { page } from '$app/state';
    import {Heading, P} from "flowbite-svelte";
    import {PUBLIC_BACKEND_URL} from "$env/static/public";

</script>

<div class="p-8 bg-white dark:bg-gray-900">
    <div class="flex-grow items-center justify-center ">
        <Heading tag="h1" class="mb-4 text-center" customSize="text-4xl font-extrabold  md:text-5xl lg:text-6xl">We've Encountered an Error</Heading>
        <P class="mb-6 text-lg lg:text-xl sm:px-16 xl:px-48 dark:text-gray-400 text-center">{page.status}: {page.error?.message}</P>
        {#if (PUBLIC_BACKEND_URL.includes("localhost"))}
            <P class="mb-6 text-lg lg:text-xl sm:px-16 xl:px-48 dark:text-gray-400 text-center">Is the Backend Running?</P>
        {/if}
    </div>
</div>
