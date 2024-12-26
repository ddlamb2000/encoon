import { env } from "$env/dynamic/private";
import { kafka } from '$lib/kafka';

export async function GET() {
  const ac = new AbortController()
  const stream = new ReadableStream({
    start(controller) {
      console.log("start controller")

      const topic = env.TOPIC_PREFIX + '-master-responses'
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

      consumer.connect(),
      consumer.subscribe({ topics: [topic] })

      try {
        consumer.run({
          eachMessage: async ({ topic, partition, message, heartbeat, pause }) => {
            const received = {
              topic: topic,
              headers: message.headers,
              key: message.key.toString(),
              value: message.value.toString()
            }
            console.log("Received from Kafka", received)
            controller.enqueue(JSON.stringify(received))
          },
        })
      } catch (error) {
        console.error(`Error subscribe to Kafka:`, error);
        return new Response(JSON.stringify({ error }), {
          headers: {
            'Content-Type': 'text/event-stream'
          },
          status: 500
        })
      }
    },
    cancel() {
      console.log("cancel and abort")
      consumer.stop();
      consumer.disconnect();
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
