import { env } from "$env/dynamic/private"
import { kafka } from '$lib/kafka'
import { newUuid } from "$lib/utils.svelte"
	
export const GET = async ({ params, request, url, cookies }) => {
  const topic = env.TOPIC_PREFIX + "-" + params.dbname + "-responses"
  const ac = new AbortController()
  const groupId = env.KAFKA_GROUP_ID + "-" + newUuid()
  cookies.set('topic', topic, { path: '/' })
  cookies.set('groupId', groupId, { path: '/' })
  console.log(`GET ${url}: Create consumer using groupId ${groupId}`)
  const consumer = kafka.consumer({
    groupId: groupId,
    minBytes: 20,
    maxBytes: 10240,
    maxWaitTimeInMs: 100,
    maxInFlightRequests: 50,
    retry: {
      retries: 5
    }
  })
  console.log(`GET stream: Start stream`)
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
                    const received = {
                      topic: batch.topic,
                      partition: batch.partition,
                      highWatermark: batch.highWatermark,
                      offset: message.offset,
                      headers: message.headers,
                      key: message.key.toString(),
                      value: message.value.toString()
                    }
                    console.log(`GET ${url}: `, received)
                    controller.enqueue(JSON.stringify(received))
                  }
                  resolveOffset(message.offset)
                  await heartbeat()
                }
            },
        })
      } catch (error) {
        console.error(`GET ${url}: Error subscribe to Kafka:`, error)
        return new Response(JSON.stringify({ error }), {
          headers: { 'Content-Type': 'text/event-stream' },
          status: 500
        })
      }
    },
    cancel() {
      console.log(`GET ${url}: Abort`)
      consumer.stop()
      consumer.disconnect()
      ac.abort()
    },
  })
  return new Response(stream, {
    headers: { 'Content-Type': 'text/event-stream' },
    status: 200
  })
}
