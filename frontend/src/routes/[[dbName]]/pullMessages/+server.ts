import { env } from "$env/dynamic/private"
import { kafka } from '$lib/kafka'
import { newUuid } from "$lib/utils.svelte.ts"
import * as metadata from "$lib/metadata.svelte"
	
export const GET = async ({ params, request, url }) => {
  if(params.dbName === undefined) {
    console.error(`PULL: Missing dbName`)
    return new Response(JSON.stringify({ error: 'Missing dbName' }), { status: 500 })
  }
  const topic = env.TOPIC_PREFIX + "-" + params.dbName + "-responses"
  const ac = new AbortController()
  const groupId = env.KAFKA_GROUP_ID + "-" + params.dbName + "-" + newUuid()
  console.log(`GET ${url}: Create consumer using groupId ${groupId}`)
  const consumer = kafka.consumer({
    groupId: groupId,
    minBytes: 20,
    maxBytes: 10240,
    maxWaitTimeInMs: 100,
    maxInFlightRequests: 50,
    retry: { retries: 5 }
  })
  console.log(`PULL: Start streaming`)
  const stream = new ReadableStream({
    start(controller) {
      try {
        console.log(`PULL: Submit an initialization message to the stream`)
        const initializationMessage = { key: metadata.InitializationKey }
        controller.enqueue(JSON.stringify(initializationMessage))
        console.log(`PULL: Start connection for consumer`)
        consumer.connect()
        console.log(`PULL: Consumer connected`)
        consumer.subscribe({ topics: [topic] })
        console.log(`PULL: Consumer subscribed to ${topic}`)
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
                console.log(`PULL: Consumer running from ${topic}`)
                for (let message of batch.messages) {
                  if (!isRunning() || isStale()) {
                    console.log(`PULL: Consumer stopping from ${topic}`)
                    break
                  }
                  if(message.key !== null && message.value !== null) {
                    const valueString = message.value.toString()
                    const received = {
                      topic: batch.topic,
                      headers: message.headers,
                      key: message.key.toString(),
                      value: valueString
                    }
                    console.log(`PULL (${valueString.length} bytes), key: ${message.key.toString()}`)
                    controller.enqueue(JSON.stringify(received))
                  } else {
                    console.log(`PULL: Message with no key nor value`, message)
                  }
                  console.log(`PULL: Resolve offset from ${topic}: ${message.offset}`)
                  resolveOffset(message.offset)
                  console.log(`PULL: Awaiting heartbeat from ${topic}`)
                  await heartbeat()
                }
            },
        })
      } catch (error) {
        console.error(`PULL: Error subscribe to Kafka:`, error)
        return new Response(JSON.stringify({ error }), { status: 500 })
      }
    },
    cancel() {
      console.log(`PULL: Abort`)
      consumer.stop()
      consumer.disconnect()
      ac.abort()
    },
  })
  console.log(`PULL: Return response as text/event-stream`)
  return new Response(stream, { headers: { 'Content-Type': 'text/event-stream' }, status: 200 })
}
