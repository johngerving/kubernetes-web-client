<script lang="ts">
    import * as Table from "$lib/components/ui/table"
    import { Skeleton } from "$lib/components/ui/skeleton/index"
	import type { PageData } from "./$types";
	import CreateWorkspaceDialog from "./CreateWorkspaceDialog.svelte";
    
    let { data, form }: { data: PageData, form: PostWorkspaceFormData } = $props();

    let workspaces: Promise<Workspace[]> = $derived(data.workspaces);
</script>

<div class="w-full h-full px-56 py-8">
    <div class="flex justify-between">
        <h1 class="text-4xl mb-3">Workspaces</h1>
        <CreateWorkspaceDialog {form}/>
    </div>
    
        {#await workspaces}
        <div class="border rounded-lg overflow-hidden">
            <Table.Root>
                <Table.Header class="bg-slate-50">
                <Table.Row>
                    <Table.Head class="w-1/3">Name</Table.Head>
                    <Table.Head class="text-center">Last Used</Table.Head>
                    <Table.Head class="text-right">Status</Table.Head>
                </Table.Row>
                </Table.Header>
                <Table.Body>
                    {#each {length: 4} as _}
                    <Table.Row>
                        <Table.Cell class="w-1/3">
                            <Skeleton class="w-full h-6 rounded-full"/>
                        </Table.Cell>
                        <Table.Cell class="text-center">
                            <Skeleton class="w-20 h-6 rounded-full m-auto"/>
                        </Table.Cell>
                        <Table.Cell class="text-right">
                            <Skeleton class="w-full h-6 rounded-full"/>
                        </Table.Cell>
                    </Table.Row>
                    {/each}
                </Table.Body>
            </Table.Root>
        </div>
        {:then workspaces}
        <div class="border rounded-lg overflow-hidden">
            <Table.Root>
                <Table.Header class="bg-slate-50">
                <Table.Row>
                    <Table.Head class="text-left w-1/3">Name</Table.Head>
                    <Table.Head class="text-center">Last Used</Table.Head>
                    <Table.Head class="text-right">Status</Table.Head>
                </Table.Row>
                </Table.Header>
                <Table.Body>
                    {#each workspaces as workspace (workspace.id)}
                        <Table.Row>
                            <Table.Cell class="text-left w-1/3">{workspace.name}</Table.Cell>
                            <Table.Cell class="text-center"></Table.Cell>
                            <Table.Cell class="text-right"></Table.Cell>
                        </Table.Row>
                    {/each}
                </Table.Body>
            </Table.Root>
        </div>
        {:catch error}
            <p class="text-red-600">Unable to retrieve workspaces</p>
        {/await}
    
</div>

