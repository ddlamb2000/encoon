import { env } from "$env/dynamic/private";
import { Kafka } from "kafkajs";

export const kafka = new Kafka({
	clientId: env.KAFKA_CLIENT_ID,
	brokers: env.KAFKA_BROKERS?.split(',') ?? [],
	retry: {
		initialRetryTime: 100,
		retries: 8
	}
})