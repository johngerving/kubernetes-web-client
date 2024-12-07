import { env } from "$env/dynamic/public"
import { fail } from "@sveltejs/kit";
import type { PageServerLoad } from "./$types.js";
import { getUserWorkspaces } from "$lib/server/utils.js";

export const load = async ({ fetch }) => {
    return {
        workspaces: getUserWorkspaces(fetch),
    }
}

export const actions = {
    create: async ({ fetch, request }) => {
        const data = await request.formData()
        
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

        if(res.status == 400) {
            if(body.message) {
                throw new Error(`Response status: ${res.status}\nMessage: ${body.message}`)
            }

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