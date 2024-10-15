<script>
    import {
        env
    } from '$env/dynamic/public';
</script>

<main>
    <div class="w-screen h-screen flex items-center justify-center">
            <a href={`${env.PUBLIC_API_URL}/auth`} class="bg-blue-500 text-xl py-3 px-5 rounded-lg text-white shadow-md hover:bg-blue-600 transition-colors duration-300">Sign in</a>
    </div>
</main>
