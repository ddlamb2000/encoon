import { d as private_env } from "./shared-server.js";
import { Kafka, CompressionTypes } from "kafkajs";
import { j as json } from "./index.js";
const kafka = new Kafka({
  clientId: private_env.KAFKA_CLIENT_ID,
  brokers: private_env.KAFKA_BROKERS?.split(",") ?? [],
  retry: {
    initialRetryTime: 100,
    retries: 8
  }
});
async function postMessage(params, request, url) {
  const topic = private_env.TOPIC_PREFIX + "-" + params.dbname + "-requests";
  const producer = kafka.producer({
    maxInFlightRequests: 50,
    allowAutoTopicCreation: true,
    retry: { retries: 5 }
  });
  try {
    await producer.connect();
  } catch (error) {
    console.error("Error connecting to Kafka:", error);
    return json({ error: "Failed to connect to Kafka." }, { status: 500 });
  }
  try {
    const data = await request.json();
    if (!data.message.trim() || data.headers.some((h) => !h.key.trim())) {
      await producer.disconnect();
      return json({ error: "Message or headers are invalid" }, { status: 400 });
    }
    await producer.send({
      topic,
      compression: CompressionTypes.GZIP,
      messages: getMessages(data),
      acks: -1
    });
    console.log(`POST ${url} to ${topic}: `, data);
    return json({ message: "Message" });
  } catch (error) {
    console.error(`Error sending message to Kafka topic ${topic}:`, error);
    return json({ error: "Failed to send message" }, { status: 500 });
  } finally {
    await producer.disconnect();
  }
}
function getMessages(req) {
  const messageKey = req.messageKey?.trim();
  const headers = Object.fromEntries(req.headers.map((h) => [h.key + "", h.value + ""]));
  if (req.selectedPartitions && req.selectedPartitions.length > 0) {
    return req.selectedPartitions.map((pId) => ({ partition: pId, key: messageKey || null, value: req.message, headers }));
  } else {
    return [{ key: messageKey || null, value: req.message, headers }];
  }
}
export {
  kafka as k,
  postMessage as p
};
