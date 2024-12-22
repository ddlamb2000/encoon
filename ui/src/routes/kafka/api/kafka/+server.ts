// src/routes/api/partitions/[topic].js
import { admin, consumer } from '$lib/kafka';
import { env } from '$env/dynamic/private';
import type { TopicInfo } from '$lib/types';

export async function GET() {

	console.log('Kafka GET');

	try {
		await consumer.connect();
	} catch (error) {
		console.error('Error connecting to Kafka:', error);
		return new Response(JSON.stringify({ error: 'Failed to connect to Kafka.' }), {
			headers: {
				'Content-Type': 'application/json'
			},
			status: 500
		});
	}

	const topic = env.TOPIC_PREFIX + '-master-responses'

	try {
		await consumer.subscribe({
			topics: [topic]
		});

		console.log(`Kafka consumer.subscribe`);

		await consumer.run({
			eachMessage: async ({ topic, partition, message, heartbeat, pause }) => {
				console.log({
					key: message.key.toString(),
					value: message.value.toString(),
					headers: message.headers,
				})
			},
		})

		return new Response("OK", {
			headers: {
				'Content-Type': 'application/json'
			},
			status: 200
		});
		
	} catch (error) {

		console.error(`Error subscribe`, error);

		return new Response(JSON.stringify({ error }), {
			headers: {
				'Content-Type': 'application/json'
			},
			status: 500
		});
	}
}
