<script lang="ts">
	import { enhance } from '$app/forms';
	import { invalidate } from '$app/navigation';
	import { Button, buttonVariants } from '$lib/components/ui/button';
    import * as Dialog from '$lib/components/ui/dialog/index'
	import { Input } from '$lib/components/ui/input';
	import { Label } from '$lib/components/ui/label';

    let { form }: { form: PostWorkspaceFormData } = $props();

    let dialogOpen = $state(false);

    // When the dialog is closed, clear the form data and errors
    $effect(() => {
        if(!dialogOpen && form != null) {
            form.errors = {};
            form.name = "";
        }
    })
</script>

<Dialog.Root bind:open={dialogOpen}>
    <Dialog.Trigger class={buttonVariants({ variant: "outline"})}>Create Workspace</Dialog.Trigger>
    <Dialog.Content>
        <Dialog.Header>
            <Dialog.Title>Create Workspace</Dialog.Title>
        </Dialog.Header>
        <form method="POST" action="?/create" use:enhance={() => {
            // Close the dialog when the form is completed
            return async ({ result, update }) => {
                await update();
                if(result.type === 'success')
                    dialogOpen = false;
            }
        }}>
            <Label for="name">Name</Label>
            <Input name="name" value={form?.name ?? ''} autocomplete="off"/>
            <!-- Display error if field is invalid -->
            {#if form?.errors?.name ?? false}
                <p class="text-red-600">{form?.errors?.name}</p>
            {/if}
            <Dialog.Footer>
                <Button type="submit" aria-label="Create" class={`${buttonVariants({ variant: "default"})} mt-4`}>Create</Button>
            </Dialog.Footer>
        </form>
    </Dialog.Content>
</Dialog.Root>