import { env } from "$env/dynamic/private"
import { kafka } from '$lib/kafka'
import { newUuid } from "$lib/utils.svelte"
	
export const GET = async ({ params, request, url }) => {
  if(params.dbName === undefined) {
    console.error(`GET ${url}: Missing dbName`)
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
  console.log(`GET ${url}: Start streaming`)
  const stream = new ReadableStream({
    start(controller) {
      try {
        console.log(`GET ${url}: Submit an initialization message to the stream`)
        const initializationMessage = {
          topic: '',
          headers: [],
          key: 'INIT',
          value: JSON.stringify({action: 'INIT'})
        }
        controller.enqueue(JSON.stringify(initializationMessage))
        console.log(`GET ${url}: Start connection for consumer`)
        consumer.connect()
        console.log(`GET ${url}: Consumer connected`)
        consumer.subscribe({ topics: [topic] })
        console.log(`GET ${url}: Consumer subscribed to ${topic}`)
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
                console.log(`GET ${url}: Consumer running from ${topic}`)
                for (let message of batch.messages) {
                  if (!isRunning() || isStale()) {
                    console.log(`GET ${url}: Consumer stopping from ${topic}`)
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
                    console.log(`GET ${url}: ${valueString.length} bytes`, received)
                    controller.enqueue(JSON.stringify(received))
                  } else {
                    console.log(`GET ${url}: Message with no key nor value`, message)
                  }
                  console.log(`GET ${url}: Resolve offset from ${topic}: ${message.offset}`)
                  resolveOffset(message.offset)
                  console.log(`GET ${url}: Awaiting heartbeat from ${topic}`)
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
  console.log(`GET ${url}: Return response as text/event-stream`)
  return new Response(stream, { headers: { 'Content-Type': 'text/event-stream' }, status: 200 })
}
