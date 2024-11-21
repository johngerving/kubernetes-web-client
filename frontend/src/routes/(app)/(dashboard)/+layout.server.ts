import { env } from "$env/dynamic/public"

/** @type {import('./$types').LayoutServerLoad} */
export async function load({ fetch }) {
    const response = await fetch(`${env.PUBLIC_API_CLUSTER_URL}/user/workspaces`)

    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`)
    }

    const json = await response.json();

    let workspaces = [];
    if(json != null)
        workspaces = json;

    return {
        workspaces: workspaces,
    }
}