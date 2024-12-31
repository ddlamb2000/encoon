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

export interface GridType {

}

export interface ColumnType {
	uuid: string
	orderNumber: number
	owned: boolean
	label: string
	name: string
	type: string
	typeUuid: string
	gridUuid: string
	grid: GridType
	gridPromptUuid: string
	bidirectional: boolean
}

export interface RowType {

}

export interface GridPost {
  rowsAdded?: RowType[]
  rowsEdited?: RowType[]
}

export interface RequestContent {
	action: string
  gridUuid?: string
	columnUuid?: string
  rowUuid?: string
	uuid?: string
	userid?: string
	password?: string
  dataSet?: GridPost
}