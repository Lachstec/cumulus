<script>
  import {
    Footer,
    FooterLinkGroup,
    FooterLink,
    FooterCopyright,
    FooterIcon,
    DarkMode,
  } from "flowbite-svelte";
  import { GithubSolid } from "flowbite-svelte-icons";
  let btnClass =
    "w-5 h-5 text-gray-500 dark:text-gray-500 hover:text-gray-900 dark:hover:text-white";
</script>

<Footer
  class="w-full p-4 bg-white border-t border-gray-200 shadow md:flex md:items-center md:justify-between md:p-6 dark:bg-gray-800 dark:border-gray-600">
  <FooterCopyright href="/" by="Bogus Beep" year={2024} />
  <FooterLinkGroup
    ulClass="flex flex-wrap items-center mt-3 text-sm text-gray-500 dark:text-gray-400 sm:mt-0">
    <FooterLink href="/about">About</FooterLink>
    <FooterLink href="/">Privacy Policy</FooterLink>
    <FooterLink href="/">Contact</FooterLink>
    <FooterIcon href="https://github.com/Lachstec/mc-hosting">
      <GithubSolid
        class="w-5 h-5 text-gray-500 dark:text-gray-500 hover:text-gray-900 dark:hover:text-white" />
    </FooterIcon>
    <FooterIcon>
      <DarkMode {btnClass} />
    </FooterIcon>
  </FooterLinkGroup>
</Footer>
