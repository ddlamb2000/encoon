import { env } from "$env/dynamic/private";
import { json } from '@sveltejs/kit'
import { Kafka, type Message, CompressionTypes } from 'kafkajs'
import { type KafkaMessageRequest, type KafkaMessageResponse } from '$lib/dataTypes.ts'

function hashCode(str) {
  var hash = 0
  for(var i = 0; i < str.length; i++) {
      hash = ~~(((hash << 5) - hash) + str.charCodeAt(i))
  }
  return hash > 0 ? hash : -hash
}

const CustomRoundRobin = () => {
  return ({ topic, partitionMetadata, message }) => {
  const nbPartitions = partitionMetadata.length
  const gridUuid = message && message.headers && message.headers.gridUuid ? message.headers.gridUuid : ""
  const hash =  hashCode(gridUuid)
  const balance = hash % nbPartitions
      return balance
  }
}

export const kafka = new Kafka({
  clientId: env.KAFKA_CLIENT_ID,
  brokers: env.KAFKA_BROKERS?.split(',') ?? [],
  retry: {
    initialRetryTime: 100,
    retries: 8
  }
})

const producer = kafka.producer({
  maxInFlightRequests: 50,
  allowAutoTopicCreation: true,
  retry: { retries: 5 },
  createPartitioner: CustomRoundRobin
})

let producerConnected = false
try {
  producer.connect()
  producerConnected = true
} catch (error) {
  console.error('Error connecting to Kafka:', error)
}


export const postMessage = async (params, request, url: string) => {
  if(params.dbName === undefined) {
    console.error('Missing dbName')
    return json({ error: 'Missing dbName' } as KafkaMessageResponse, { status: 500 })
  }
  const topic = env.TOPIC_PREFIX + "-" + params.dbName + "-requests"
  if(!producerConnected) {
    return json({ error: 'Failed to connect to Kafka' } as KafkaMessageResponse, { status: 500 })
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
    return json({  } as KafkaMessageResponse)
  } catch (error) {
    console.error(`Error sending message to Kafka topic ${topic}:`, error)
    return json({ error: 'Failed to send message' } as KafkaMessageResponse, { status: 500 })
  }
}

const getMessages = (req: KafkaMessageRequest): Message[] => {
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