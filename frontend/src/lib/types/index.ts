export interface KafkaMessageHeader {
	key: string
	value: string
}

export interface KafkaMessageRequest {
	messageKey?: string
	message: string
	headers: KafkaMessageHeader[]
	selectedPartitions: number[]
}

export interface KafkaMessageResponse {
	message: string
	error?: string
}

export type Partition = number