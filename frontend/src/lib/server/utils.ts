import { env } from "$env/dynamic/public";

const delay = (delayInms: number) => {
    return new Promise(resolve => setTimeout(resolve, delayInms));
};

export async function getUserWorkspaces(fetch: (input: RequestInfo | URL, init?: RequestInit) => Promise<Response>): Promise<Workspace[]> {
    const response = await fetch(`${env.PUBLIC_API_CLUSTER_URL}/user/workspaces`)

    if (!response.ok) {
        throw new Error(`Response status: ${response.status}`)
    }

    const json = await response.json();

    let workspaces: Workspace[] = [];
    if(json != null)
        workspaces = json;

    await delay(2000);

    return workspaces;
}