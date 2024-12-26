import { env } from "$env/dynamic/private";
import { kafka } from '$lib/kafka';

export async function GET() {
  const topic = env.TOPIC_PREFIX + '-master-responses'
  const ac = new AbortController()
  console.log("GET Kafak stream: Start stream")
  const stream = new ReadableStream({
    start(controller) {
      const consumer = kafka.consumer({
        groupId: env.KAFKA_GROUP_ID,
        minBytes: 20,
        maxBytes: 1024,
        maxWaitTimeInMs: 10,
        maxInFlightRequests: 50,
        retry: {
          retries: 5
        }
      })

      try {
        console.log("GET Kafak stream: Connect Kafka consumer")
        consumer.connect()
        console.log(`GET Kafak stream: Subscribe Kafka consumer to ${topic}`)
        consumer.subscribe({ topics: [topic] })
        consumer.run({
          eachMessage: async ({ topic, partition, message, heartbeat, pause }) => {
            const received = {
              topic: topic,
              headers: message.headers,
              key: message.key.toString(),
              value: message.value.toString()
            }
            console.log("GET Kafak stream: Received from Kafka", received)
            controller.enqueue(JSON.stringify(received))
          },
        })
      } catch (error) {
        console.error(`GET Kafak stream: Error subscribe to Kafka:`, error);
        return new Response(JSON.stringify({ error }), {
          headers: {
            'Content-Type': 'text/event-stream'
          },
          status: 500
        })
      }
    },
    cancel() {
      console.log("GET Kafak stream: Stop Kafka consumer")
      consumer.stop()
      console.log("GET Kafak stream: Disconnect Kafka consumer")
      consumer.disconnect()
      console.log("GET Kafak stream: Cancel stream")
      stream.cancel()
      console.log("GET Kafak stream: Abort")
      ac.abort()
    },
  })
  return new Response(stream, {
    headers: {
      'Content-Type': 'text/event-stream'
    },
    status: 200
  })
}
