import { env } from "$env/dynamic/private";
import { Kafka } from "kafkajs";

const kafka = new Kafka({
	clientId: env.KAFKA_CLIENT_ID,
	brokers: env.KAFKA_BROKERS?.split(',') ?? []
});

export const producer = kafka.producer({
	maxInFlightRequests: 50,
	retry: {
		retries: 5,
		initialRetryTime: 100
	}
});

export const consumer = kafka.consumer({
	groupId: env.KAFKA_GROUP_ID
});

export function connect() {
	return Promise.all([
		producer.connect(),
		consumer.connect()
	]);
}