import { env } from "$env/dynamic/private";
import { Kafka } from "kafkajs";
import { json } from '@sveltejs/kit'
import { type Message, CompressionTypes } from 'kafkajs'
import type { KafkaMessageRequest, KafkaMessageResponse } from '$lib/types'

export const kafka = new Kafka({
	clientId: env.KAFKA_CLIENT_ID,
	brokers: env.KAFKA_BROKERS?.split(',') ?? [],
	retry: {
		initialRetryTime: 100,
		retries: 8
	}
})

export async function postMessage(params, request, url) {
	const topic = env.TOPIC_PREFIX + "-" + params.dbname + "-requests"
	const producer = kafka.producer({
		maxInFlightRequests: 50,
		allowAutoTopicCreation: true,
		retry: { retries: 5 }
	})
	try {
		await producer.connect()
	} catch (error) {
		console.error('Error connecting to Kafka:', error)
		return json({ error: 'Failed to connect to Kafka.' } as KafkaMessageResponse, { status: 500 })
	}
	try {
		const data: KafkaMessageRequest = await request.json()
		if (!data.message.trim() || data.headers.some((h) => !h.key.trim())) {
			await producer.disconnect()
			return json({ error: 'Message or headers are invalid' }, { status: 400 })
		}
		await producer.send({
			topic: topic,
			compression: CompressionTypes.GZIP,
			messages: getMessages(data),
			acks: -1
		})
		const dataLength = JSON.stringify(data).length
		console.log(`POST ${url} to ${topic}: ${dataLength} bytes`, data)
		return json({ message: 'Message sent successfully' } as KafkaMessageResponse)
	} catch (error) {
		console.error(`Error sending message to Kafka topic ${topic}:`, error)
		return json({ error: 'Failed to send message' } as KafkaMessageResponse, { status: 500 })
	} finally { 
		await producer.disconnect()
	}
}

function getMessages(req: KafkaMessageRequest): Message[] {
	const messageKey = req.messageKey?.trim()
	const headers = Object.fromEntries(req.headers.map((h) => [ h.key + "", h.value + ""]))
	if(req.selectedPartitions && req.selectedPartitions.length > 0) {
		return req.selectedPartitions.map((pId) => (
			{ partition: pId, key: messageKey || null, value: req.message, headers: headers }
		))
	} else {
		return [ { key: messageKey || null, value: req.message, headers: headers } as Message ]
	}
}