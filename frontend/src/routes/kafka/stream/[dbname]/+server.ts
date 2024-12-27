import { env } from "$env/dynamic/private";
import { kafka } from '$lib/kafka';

export const GET = async ({ params, request, url, cookies }) => {
  
  const topic = env.TOPIC_PREFIX + "-" + params.dbname + "-responses"
  const ac = new AbortController()
  const groupId = env.KAFKA_GROUP_ID
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
        console.log("GET Kafak stream: Connect Kafka consumer")
        consumer.connect()
        console.log(`GET Kafak stream: Subscribe Kafka consumer to ${topic}`)
        consumer.subscribe({ topics: [topic] })
        consumer.run({
          eachMessage: async ({ topic, partition, message, heartbeat, pause }) => {
            if(message.key !== null && message.value !== null) {
              const received = {
                topic: topic,
                headers: message.headers,
                key: message.key.toString(),
                value: message.value.toString()
              }
              console.log("GET Kafak stream: Received from Kafka", received)
              controller.enqueue(JSON.stringify(received))
            }
            await heartbeat()
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
