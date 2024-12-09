import { env } from "$env/dynamic/public"
import { fail } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types.js";
import { getUserWorkspaces } from "$lib/server/utils.js";

export const load: PageServerLoad = async ({ fetch }) => {
    const workspaces = getUserWorkspaces(fetch);
    workspaces.catch((e) => console.log(e)); // Catch a rejected promise

    return {
        workspaces: workspaces,
    }
}

export const actions = {
    create: async ({ fetch, request }) => {
        const data = await request.formData()
        
        // Make a POST request to create a new workspace
        const res = await fetch(
            `${env.PUBLIC_API_CLUSTER_URL}/user/workspaces`,
            {
                method: "POST",
                headers: {
                    "content-type": "application/json"
                },
                body: JSON.stringify({
                    "name": data.get("name"),
                })
            }
        )

        const body = await res.json()

        // Throw an error if the response was unsuccessful
        if(res.status == 400) {
            if(body.message) {
                throw new Error(`Response status: ${res.status}\nMessage: ${body.message}`)
            }

            // Send errors and form data to client
            const errors: PostWorkspaceFormErrors = {
                name: body.name
            };

            return fail(400, {
                name: data.get("name"),
                errors: errors
            })
        }
    }
}