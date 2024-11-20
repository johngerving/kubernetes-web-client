import { env } from "$env/dynamic/public"

/** @type {import('@sveltejs/kit').HandleFetch} */
export async function handleFetch({ event, request, fetch }) {
    // Redirect fetch if the request is to the backend URL
    if(request.url.startsWith(env.PUBLIC_API_URL)) {
        request = new Request(
            request.url,
            request,
        )

        // Copy the cookies to a new request
        request.headers.set(
            'cookie',
            event.cookies
                .getAll()
                .filter(({ value }) => value !== '') // Account for cookie deleted in current request
                .map(({ name, value }) => `${name}=${encodeURIComponent(value)}`)
                .join('; ')
        );

    }
    return fetch(request);
}