// src/routes/api/partitions/[topic].js
import { consumer } from '$lib/kafka';
import { env } from '$env/dynamic/private';

const topic = env.TOPIC_PREFIX + '-master-responses'
export async function GET() {
  const ac = new AbortController()
  const stream = new ReadableStream({
    start(controller) {
      try {
        consumer.subscribe({ topics: [topic] })
        console.log(`Kafka consumer subscribed to ${topic}`);
        consumer.run({
          eachMessage: async ({ topic, partition, message, heartbeat, pause }) => {
            console.log(topic,{
              headers: message.headers,
              key: message.key.toString(),
              value: message.value.toString()
            })
            controller.enqueue(message.value)
          },
        })
      } catch (error) {
        console.error(`Error subscribe to Kafka:`, error);
      }
    },
    cancel() {
      console.log("cancel and abort")
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
