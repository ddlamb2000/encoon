import type { RequestHandler } from '@sveltejs/kit'
import { postMessage } from '$lib/kafka'

export const POST: RequestHandler = async ({ params, request }) => {
    return postMessage(params, request)
}