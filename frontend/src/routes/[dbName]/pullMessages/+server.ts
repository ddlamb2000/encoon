import { env } from "$env/dynamic/private"
import { kafka } from '$lib/kafka'
import { newUuid } from "$lib/utils.svelte"
	
export const GET = async ({ params, request, url }) => {
  if(params.dbName === undefined) {
    console.error('Missing dbName')
    return new Response(JSON.stringify({ error: 'Missing dbName' }), { status: 500 })
  }
  const topic = env.TOPIC_PREFIX + "-" + params.dbName + "-responses"
  const ac = new AbortController()
  const groupId = env.KAFKA_GROUP_ID + "-" + newUuid()
  console.log(`GET ${url}: Create consumer using groupId ${groupId}`)
  const consumer = kafka.consumer({
    groupId: groupId,
    minBytes: 20,
    maxBytes: 10240,
    maxWaitTimeInMs: 100,
    maxInFlightRequests: 50,
    retry: { retries: 5 }
  })
  const stream = new ReadableStream({
    start(controller) {
      try {
        console.log(`GET ${url}: Subscribe consumer to ${topic}`)
        consumer.connect()
        consumer.subscribe({ topics: [topic] })
        consumer.run({
            eachBatchAutoResolve: true,
            autoCommitInterval: 50,
            eachBatch: async ({
                batch,
                resolveOffset,
                heartbeat,
                commitOffsetsIfNecessary,
                uncommittedOffsets,
                isRunning,
                isStale,
                pause,
            }) => {
                for (let message of batch.messages) {
                  if (!isRunning() || isStale()) break
                  if(message.key !== null && message.value !== null) {
                    const valueString = message.value.toString()
                    const received = {
                      topic: batch.topic,
                      headers: message.headers,
                      key: message.key.toString(),
                      value: valueString
                    }
                    console.log(`GET ${url}: ${valueString.length} bytes`, received)
                    controller.enqueue(JSON.stringify(received))
                  }
                  resolveOffset(message.offset)
                  await heartbeat()
                }
            },
        })
      } catch (error) {
        console.error(`GET ${url}: Error subscribe to Kafka:`, error)
        return new Response(JSON.stringify({ error }), { status: 500 })
      }
    },
    cancel() {
      console.log(`GET ${url}: Abort`)
      consumer.stop()
      consumer.disconnect()
      ac.abort()
    },
  })
  return new Response(stream, { headers: { 'Content-Type': 'text/event-stream' }, status: 200 })
}
