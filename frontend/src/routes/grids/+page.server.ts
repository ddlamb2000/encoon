import type { PageServerLoad } from './$types';
import { newUuid } from "$lib/utils.svelte"
import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types';

const delay = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export const load: PageServerLoad = async () => {
  const getQuote = async () => {
    await delay(2000);
    return "The only way to do great work is to love what you do. - Steve Jobs";
  }
  return {
    quote: getQuote()
  }
}
