<script lang="ts">
    import {onMount} from "svelte";
    import {
        Button,
        Table,
        TableBody,
        TableBodyCell,
        TableBodyRow,
        TableHead,
        TableHeadCell
    } from "flowbite-svelte";
    import auth from "$lib/service/auth_service";
    import type { Auth0Client } from "@auth0/auth0-spa-js";
    
    type UserData = {
        ID: number,
        Sub: "",
        name: string,
        role: string,
    }

    let data: UserData[]
    let auth0Client: Auth0Client;

    async function getAllUsers() {
        auth0Client = await auth.createClient();
        const token = await auth0Client.getTokenSilently();
        const res = await fetch("http://localhost:10000/users", {
            method: "GET",
            headers: {
                Authorization: `Bearer ${token}`
            }
        });
        data = await res.json();
    }

    onMount(async () => {
        await getAllUsers();
    });

    let isVisible = false;

    function deleteUser(ID : number) {
        fetch(`http://localhost:10000/users/${ID}`, {
            method: "DELETE",
        }).then(res => {
            if(res.status === 204){
                isVisible = true;
                getAllUsers();
            } else {
                console.error('Failed to delete user: ', res.status);
            }
        })
        .catch(error => {
            console.error(error);
        });
    }
</script>

<ul>
    <div class="p-8 mt-16 bg-white dark:bg-gray-900 h-screen">
        <Table>
            <caption
                    class="p-5 text-lg font-semibold text-left text-gray-900 bg-white dark:text-white dark:bg-gray-800">
                Users
                <p class="mt-1 text-sm font-normal text-gray-500 dark:text-gray-400">
                    View, Edit and Delete Users
                </p>
            </caption>
            <TableHead>
                <TableHeadCell>Username</TableHeadCell>
                <TableHeadCell>Role</TableHeadCell>
                <TableHeadCell>Options</TableHeadCell>
            </TableHead>
            <TableBody tableBodyClass="divide-y">
                {#each data as { name, role, ID }}
                    <TableBodyRow>
                        <TableBodyCell>{name}</TableBodyCell>
                        <TableBodyCell>{role}</TableBodyCell>
                        <TableBodyCell>
                            <Button size="xs" outline color="green"
                                href="../../user/{ID}"
                            >View</Button>
                            <Button size="xs" outline color="yellow"

                            >Edit</Button>
                            <Button size="xs" outline color="red"
                                role="button"
                                on:click={() => {
                                        if(confirm("Do you really want to delete this user?")){
                                            deleteUser(ID)
                                            console.log("Deleted user")
                                        } else {
                                            console.log('Not Deleted');
                                        }
                                    }}
                            >Delete</Button>
                        </TableBodyCell>
                    </TableBodyRow>
                {/each}
            </TableBody>
        </Table>
    </div>
</ul>