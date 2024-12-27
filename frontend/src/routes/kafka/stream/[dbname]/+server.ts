import { env } from "$env/dynamic/private"
import { kafka } from '$lib/kafka'
import { newUuid } from "$lib/utils.svelte"
	
export const GET = async ({ params, request, url, cookies }) => {
  const topic = env.TOPIC_PREFIX + "-" + params.dbname + "-responses"
  const ac = new AbortController()
  const groupId = env.KAFKA_GROUP_ID + "-" + newUuid()
  console.log(`GET Kafak stream: Create Kafka consumer using groupId ${groupId}`)
  const consumer = kafka.consumer({
    groupId: groupId,
    minBytes: 20,
    maxBytes: 1024,
    maxWaitTimeInMs: 10,
    maxInFlightRequests: 50,
    retry: {
      retries: 5
    }
  })
  console.log(`GET Kafak stream: Start stream`)
  const stream = new ReadableStream({
    start(controller) {
      try {
        console.log(`GET Kafak stream: Subscribe Kafka consumer to ${topic}`)
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
                    console.log("GET Kafak stream: Received from Kafka", received)
                    controller.enqueue(JSON.stringify(received))
                  }
                  resolveOffset(message.offset)
                  await heartbeat()
                }
            },
        })
      } catch (error) {
        console.error(`GET Kafak stream: Error subscribe to Kafka:`, error)
        return new Response(JSON.stringify({ error }), {
          headers: { 'Content-Type': 'text/event-stream' },
          status: 500
        })
      }
    },
    cancel() {
      console.log("GET Kafak stream: Abort")
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
