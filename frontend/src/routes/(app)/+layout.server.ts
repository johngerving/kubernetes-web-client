import { env } from "$env/dynamic/public"
import { redirect } from "@sveltejs/kit"

/** @type {import('./$types').LayoutServerLoad} */
export async function load({ fetch }) {
    console.log("url:", env.PUBLIC_API_CLUSTER_URL)
    const response = await fetch(`${env.PUBLIC_API_CLUSTER_URL}/user`)

    // Redirect if not logged in
    if(response.status == 401)
        redirect(302, '/login')
    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`)
    }

    const json = await response.json();

    return {
        user: json
    }
}