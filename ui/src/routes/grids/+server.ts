import { consumer } from '$lib/kafka';
import { env } from '$env/dynamic/private';

console.log('Connection to Kafka as a consumer');

try {
    await consumer.connect();
} catch (error) {
    console.error('Error connecting to Kafka:', error);
}

const topic = env.TOPIC_PREFIX + '-master-responses'

try {
    await consumer.subscribe({
        topics: [topic]
    });

    console.log(`Kafka consumer subscribed on ${topic}`);

    await consumer.run({
        eachMessage: async ({ topic, partition, message, heartbeat, pause }) => {
            console.log({
                topic: topic.toString(),
                key: message.key.toString(),
                value: message.value.toString(),
                headers: message.headers,
            })
        },
    })
    
} catch (error) {

    console.error(`Error subscribe`, error);
}
