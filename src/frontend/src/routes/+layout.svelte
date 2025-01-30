<script>
    import Footer from "$lib/components/+Footer.svelte";
    import Navbar from "$lib/components/+Navbar.svelte";
    import "../app.css";
</script>

<div class="flex flex-col min-h-screen">
    <Navbar />

    <!-- Main Content with Proper Padding -->
    <main class="flex-grow pt-[55px] pb-[55spx]">
        <slot />
    </main>

    <Footer />
</div>
