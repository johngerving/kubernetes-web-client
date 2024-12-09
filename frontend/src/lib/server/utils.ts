import { env } from "$env/dynamic/public";
import { error } from "@sveltejs/kit";

export async function getUserWorkspaces(fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>): Promise<Workspace[]|never> {
    // Use fetch function passed from form
    const response = await fetch(`${env.PUBLIC_API_CLUSTER_URL}/user/workspaces`)

    // If there was an error, return a rejected promise
    if (!response.ok) {
        const promise = Promise.reject(new Error("unable to retrieve workspaces"));
        return promise;
    }

    const json = await response.json();

    let workspaces: Workspace[] = [];
    if(json != null)
        workspaces = json;

    return workspaces;
}