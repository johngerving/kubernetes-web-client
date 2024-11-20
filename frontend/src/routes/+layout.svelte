<script lang="ts">
    import "../app.css";
    import * as Avatar from "$lib/components/ui/avatar";

    /**
     * @typedef {Object} Props
     * @property {import('svelte').Snippet} [children]
     */

    /** @type {Props} */
    /** @type {{ data: import('./$types').LayoutData, children: Snippet }}*/
    let { data, children } = $props();

    const user = data.user as User;
</script>

<nav class="flex justify-end p-2 bg-slate-100">
    <p>{user.email}</p>
    <button>
        <Avatar.Root class="">
            <Avatar.Image src="https://github.com/shadcn.png" alt="@shadcn" />
            <Avatar.Fallback>CN</Avatar.Fallback>
        </Avatar.Root>
    </button>
</nav>

<main>
    {@render children?.()}
</main>
