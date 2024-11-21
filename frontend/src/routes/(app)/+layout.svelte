<script lang="ts">
    import "../../app.css";
    import * as Avatar from "$lib/components/ui/avatar";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
    import LogOut from "lucide-svelte/icons/log-out";
	import { userInitial } from "$lib/utils";
	import { env } from "$env/dynamic/public";

    /**
     * @typedef {Object} Props
     * @property {import('svelte').Snippet} [children]
     */

    /** @type {Props} */
    /** @type {{ data: import('./$types').LayoutData, children: Snippet }}*/
    let { data, children } = $props();

    const user = data.user as User;
    
    const initial = userInitial(user);
</script>

<nav class="flex justify-end p-2 border-b border-b-gray-100">
    <DropdownMenu.Root>
        <DropdownMenu.Trigger>
            <Avatar.Root>
                <Avatar.Fallback><span class="text-slate-900 font-bold text-xl">{initial}</span></Avatar.Fallback>
            </Avatar.Root>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
            <DropdownMenu.Group>
                <DropdownMenu.Label>{user.email}</DropdownMenu.Label>
                <DropdownMenu.Separator />
                <DropdownMenu.Item>
                    <form action={`${env.PUBLIC_API_URL}/auth/logout`} method="post" class="w-full">
                        <button type="submit" class="w-full text-left">
                            <LogOut class="mr-2 h-4 w-4 inline" />
                            <span class="inline">Log out</span>
                        </button>
                    </form>
                </DropdownMenu.Item>
            </DropdownMenu.Group>
        </DropdownMenu.Content>
    </DropdownMenu.Root>
    
</nav>

<main>
    {@render children?.()}
</main>
