<script>
  import { Navbar, NavBrand, NavLi, NavUl } from "flowbite-svelte";
  import Useritem from "./+Useritem.svelte";
</script>

<div class="relative px-8 mb-16">
  <Navbar class="px-2 sm:px-4 py-2.5 fixed w-full z-20 top-0 start-0 border-b">
    <NavBrand href="/">
      <img src="/images/logo.png" class="me-3 h-6 sm:h-9" alt="Hosting Logo" />
      <span
        class="self-center whitespace-nowrap text-xl font-semibold dark:text-white"
        >MC-Hosting</span>
    </NavBrand>
    <NavUl>
      <NavLi href="/" active={true}>Home</NavLi>
      <NavLi href="/server">Server</NavLi>
      <NavLi href="/pricing">Pricing</NavLi>
      <NavLi href="/user">Users</NavLi>
    </NavUl>
    <Useritem></Useritem>
  </Navbar>
</div>
