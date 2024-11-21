<script>
    import {
        env
    } from '$env/dynamic/public';

    import "../../app.css";
</script>

<main>
    <form action={`${env.PUBLIC_API_URL}/auth/login`} method="post" class="w-screen h-screen flex items-center justify-center">
        <button type="submit" class="bg-blue-500 text-xl py-3 px-5 rounded-lg text-white shadow-md hover:bg-blue-600 transition-colors duration-300">Sign in</button>
    </form>
</main>
