<script lang="ts">
    import * as Table from "$lib/components/ui/table"
    import * as Dialog from "$lib/components/ui/dialog";
    import {
        Button,
        buttonVariants
    } from "$lib/components/ui/button/index.js";
    /** @type {Props} */
    let { data } = $props();

    const workspaces: Workspace[] = data.workspaces;
</script>

<div class="w-full h-full px-56 py-8">
    <div class="flex justify-between">
        <h1 class="text-4xl mb-3">Workspaces</h1>
        <Dialog.Root>
        <Dialog.Trigger class={buttonVariants({ variant: "outline"})}>Create Workspace</Dialog.Trigger>
        <Dialog.Content>
            <Dialog.Header>
            <Dialog.Title>Are you sure absolutely sure?</Dialog.Title>
            <Dialog.Description>
            This action cannot be undone. This will permanently delete your account
            and remove your data from our servers.
            </Dialog.Description>
            </Dialog.Header>
        </Dialog.Content>
        </Dialog.Root>
    </div>
    <div class="border rounded-lg overflow-hidden">
        <Table.Root>
            <Table.Header class="bg-slate-50">
             <Table.Row>
              <Table.Head class="">Name</Table.Head>
              <Table.Head class="text-center">Last Used</Table.Head>
              <Table.Head class="text-right">Status</Table.Head>
             </Table.Row>
            </Table.Header>
            <Table.Body>
                {#each workspaces as workspace}
                    <Table.Row>
                        <Table.Cell>{workspace.name}</Table.Cell>
                        <Table.Cell class="text-center"></Table.Cell>
                        <Table.Cell class="text-right"></Table.Cell>
                    </Table.Row>
                {/each}
            </Table.Body>
        </Table.Root>
    </div>
</div>

