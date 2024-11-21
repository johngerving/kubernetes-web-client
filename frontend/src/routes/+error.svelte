<script lang="ts">
    import { page } from "$app/stores";
</script>

<main class="w-screen h-screen flex justify-center items-center">
    <div id="error" class="text-center">
        <h1 class="text-7xl mb-3 text-slate-900">{$page.status}</h1>
        <p class="text-3xl text-slate-700">{$page.error!.message}</p>
    </div>
</main>