<script>
    import { Avatar, DropdownDivider, DropdownHeader, DropdownItem, NavHamburger, Dropdown} from "flowbite-svelte";

</script>
<div class="flex items-center space-x-4" id="user-menu">
    <Avatar id="avatar-menu" src="/images/example-profile.png" />
    <div class="space-y-1 font-medium dark:text-white">
        Login
    </div>
    <NavHamburger class1="w-full md:flex md:w-auto md:order-1" />
</div>
<Dropdown placement="bottom" triggeredBy="#user-menu">
    <DropdownHeader>
        <span class="block text-sm">Beep Boop</span>
        <span class="block truncate text-sm font-medium">name@example.com</span>
    </DropdownHeader>
    <DropdownItem>Dashboard</DropdownItem>
    <DropdownItem>Settings</DropdownItem>
    <DropdownItem href="/user/1/servers">My Servers</DropdownItem>
    <DropdownDivider />
    <DropdownItem>Sign out</DropdownItem>
</Dropdown>