<script lang="ts">
    import "../app.css";
    import * as Avatar from "$lib/components/ui/avatar";
    import * as DropdownMenu from "$lib/components/ui/dropdown-menu";
	import { userInitial } from "$lib/utils";

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

<nav class="flex justify-end p-2 bg-slate-100">
    <DropdownMenu.Root>
        <DropdownMenu.Trigger>
            <Avatar.Root>
                <Avatar.Fallback><span class="text-slate-900 font-bold text-xl">{initial}</span></Avatar.Fallback>
            </Avatar.Root>
        </DropdownMenu.Trigger>
        <DropdownMenu.Content>
            <DropdownMenu.Group>
                <DropdownMenu.Label>My Account</DropdownMenu.Label>
                <DropdownMenu.Separator />
                <DropdownMenu.Item>Profile</DropdownMenu.Item>
                <DropdownMenu.Item>Billing</DropdownMenu.Item>
                <DropdownMenu.Item>Team</DropdownMenu.Item>
                <DropdownMenu.Item>Subscription</DropdownMenu.Item>
            </DropdownMenu.Group>
        </DropdownMenu.Content>
    </DropdownMenu.Root>
    
</nav>

<main>
    {@render children?.()}
</main>
