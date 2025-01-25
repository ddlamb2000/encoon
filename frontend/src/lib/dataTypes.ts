// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

export interface KafkaMessageHeader {
  key: string
  value: string
}

export interface KafkaMessageRequest {
  messageKey?: string
  message: string
  headers: KafkaMessageHeader[]
}

export interface KafkaMessageResponse {
  message: string
  error?: string
}

export interface GridType extends RowType {
  columns?: ColumnType[]
  columnsUsage?: ColumnType[]
}

export interface ColumnType {
  uuid: string
  orderNumber?: number
  owned?: boolean
  label?: string
  name: string
  type: string
  typeUuid: string
  gridUuid: string
  grid?: GridType
  gridPromptUuid?: string
  bidirectional?: boolean
}

export interface RowType {
  gridUuid: string
	uuid: string
	text1?: string
	text2?: string
	text3?: string
	text4?: string
	text5?: string
	text6?: string
	text7?: string
	text8?: string
	text9?: string
	text10?: string
  int1?: number
  int2?: number
  int3?: number
  int4?: number
  int5?: number
  int6?: number
  int7?: number
  int8?: number
  int9?: number
  int10?: number
  displayString?: string
  references?: ReferenceType[]
  created?: Date
  updated?: Date
}

export interface ReferenceType {
	owned: boolean
	label?: string
	name?: string
	gridUuid?: string
	rows?: RowType[]
}

export interface GridPost {
  rowsAdded?: RowType[]
  rowsEdited?: RowType[]
  rowsDeleted?: RowType[]
	referencedValuesAdded?: GridReferencePost[]
	referencedValuesRemoved?: GridReferencePost[]
}

export interface GridReferencePost {
	columnName: string
	fromUuid: string
	toGridUuid: string
	uuid: string
	owned: boolean
}

export interface RequestContent {
  action: string
  actionText?: string
  gridUuid?: string
  columnUuid?: string
  uuid?: string
  userid?: string
  password?: string
  filterColumnOwned?: boolean
  filterColumnName?: string
  filterColumnGridUuid?: string
  filterColumnValue?: string
  dataSet?: GridPost
}

export interface ResponseContent {
  action: string
  actionText?: string
  status: string
  gridUuid?: string
  columnUuid?: string
  uuid?: string
	textMessage?: string
	firstName?: string
	lastName?: string
	jwt?: string
  dataSet?: GridResponse
}

export interface GridResponse {
  grid: GridType
  countRows: number
  rows: RowType[]
  rowsAdded?: RowType[]
  rowsEdited?: RowType[]
  rowsDeleted?: RowType[]
  gridUuid?: string
  columnUuid?: string
  uuid?: string
	referencedValuesAdded?: GridReferencePost[]
	referencedValuesRemoved?: GridReferencePost[]
  canViewRows: boolean
  canEditRows: boolean
  canAddRows: boolean
  canEditGrid: boolean
  filterColumnOwned?: boolean
  filterColumnName?: string
  filterColumnGridUuid?: string
  filterColumnValue?: string
}

export interface TransactionItem {
  messageKey: string,
  action: string
  actionText?: string
  status?: string
  gridUuid?: string
  dateTime?: string
  sameContext?: boolean
  timeOut?: boolean
}

export interface Transaction {
  request?: TransactionItem
  response?: TransactionItem
}