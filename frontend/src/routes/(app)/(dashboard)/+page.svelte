<script lang="ts">
    import * as Table from "$lib/components/ui/table"
    import * as Dialog from "$lib/components/ui/dialog";
    import {
        Button,
        buttonVariants
    } from "$lib/components/ui/button/index.js";
	import { enhance } from "$app/forms";
	import { FormControl, FormField } from "$lib/components/ui/form";
	import FormLabel from "$lib/components/ui/form/form-label.svelte";
	import Input from "$lib/components/ui/input/input.svelte";
	import type { PageData } from "./$types";
	import { Label } from "$lib/components/ui/label";
	import CreateWorkspaceDialog from "./CreateWorkspaceDialog.svelte";
    
    let { data, form }: { data: PageData, form: PostWorkspaceFormData } = $props();

    let workspaces: Workspace[] = $derived(data.workspaces);
</script>

<div class="w-full h-full px-56 py-8">
    <div class="flex justify-between">
        <h1 class="text-4xl mb-3">Workspaces</h1>
        <CreateWorkspaceDialog {form}/>
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
                {#each workspaces as workspace (workspace.id)}
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

