import { env } from "$env/dynamic/public"

/** @type {import('./$types').LayoutServerLoad} */
export async function load({ fetch }) {
    try {
        const response = await fetch(`http://api-svc.default.svc.cluster.local:80/user`)
        if (!response.ok) {
            throw new Error(`Response status: ${response.status}`)
        }

        const json = await response.json();

        return {
            user: json
        }
    } catch (error: any) {
        console.error(error);
    }
     
    return {
        user: null,
    }
}