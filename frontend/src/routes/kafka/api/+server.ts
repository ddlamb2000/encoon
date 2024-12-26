import type { RequestHandler } from '@sveltejs/kit';
import { json } from '@sveltejs/kit';
import { type Message, CompressionTypes } from 'kafkajs';
import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types';
import { env } from '$env/dynamic/private';
import { kafka } from '$lib/kafka';
import type { UserTextMessage } from '$lib/server';

export const POST: RequestHandler = async ({ request }) => {

	const producer = kafka.producer({
		maxInFlightRequests: 50,
		allowAutoTopicCreation: true,
		retry: {
			retries: 5
		}
	})

	try {
		await producer.connect();
	} catch (error) {
		console.error('Error connecting to Kafka:', error);
		return json({ error: 'Failed to connect to Kafka.' } as KafkaMessageResponse, { status: 500 });
	}

	const topic = env.TOPIC_PREFIX + '-master-requests'

	try {
		const data: KafkaMessageRequest = await request.json();

		// Validate the data before sending to Kafka
		if (!data.message.trim() || data.headers.some((h) => !h.key.trim())) {
			await producer.disconnect();
			return json({ error: 'Message or headers are invalid.' }, { status: 400 });
		}

		await producer.send({
			topic: topic,
			compression: CompressionTypes.GZIP,
			messages: getMessages(data),
			acks: -1
		});

		return json({ message: 'Message sent successfully.' } as KafkaMessageResponse);
	} catch (error) {
		console.error(`Error sending message to Kafka topic ${topic}:`, error);
		return json({ error: 'Failed to send message.' } as KafkaMessageResponse, { status: 500 });
	} finally { 
		await producer.disconnect();
	}
};

function getMessages(req: KafkaMessageRequest): Message[] {
	const messageKey = req.messageKey?.trim();

	const value = JSON.stringify({
			userText: req.message
		} as UserTextMessage
	);

	const headers = Object.fromEntries(req.headers.map((h) => [ h.key + "", h.value + ""]));

	if (req.selectedPartitions && req.selectedPartitions.length > 0) {
		return req.selectedPartitions.map((pId) => ({
			partition: pId,
			key: messageKey || null,
			value: value,
			headers: headers
		}));
	} else {
		return [
			{
				key: messageKey || null,
				value: value,
				headers: headers
			} as Message
		];
	}
}
