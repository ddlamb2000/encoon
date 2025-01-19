import type { RowType, ColumnType, GridType } from '$lib/dataTypes.ts'

export class Focus {
  #grid?: GridType | undefined = $state(undefined)
  #column?: ColumnType | undefined = $state(undefined)
  #row?: RowType | undefined = $state(undefined)

  reset = () => {
    this.#grid = undefined
    this.#column = undefined
    this.#row = undefined
  }

  set = (grid: GridType | undefined, column: ColumnType | undefined, row: RowType | undefined) => {
    this.#grid = grid
    this.#column = column
    this.#row = row
  }

  hasFocus = (): boolean => this.#grid !== undefined
  hasGrid = (): boolean => this.#grid !== undefined
  hasColumn = (): boolean => this.#column !== undefined
  hasRow = (): boolean => this.#row !== undefined

  isColumnFocused = (grid: GridType, column: ColumnType): boolean | undefined => {
    return this.#grid !== undefined && this.#grid.uuid === grid.uuid 
            && this.#column !== undefined && this.#column.uuid === column.uuid
  }

  isRowFocused = (grid: GridType, row: RowType): boolean | undefined => {
    return this.#grid !== undefined && this.#grid.uuid === grid.uuid 
            && this.#row !== undefined && this.#row.uuid === row.uuid 
  }

  isFocused = (grid: GridType, column: ColumnType, row: RowType): boolean | undefined => {
    return this.#grid !== undefined && this.#grid.uuid === grid.uuid 
            && this.#row !== undefined && this.#row.uuid === row.uuid 
            && this.#column !== undefined && this.#column.uuid === column.uuid
  }

  getGridName = () => this.#grid !== undefined ? this.#grid.text1 : ""
  getGridUuid = () => this.#grid !== undefined ? this.#grid.uuid : ""
  getColumnName = () => this.#column !== undefined ? this.#column.label : ""
  getColumnType = () => this.#column !== undefined ? this.#column.type : ""
  getRowName = () => this.#row !== undefined ? this.#row.displayString : ""
  getCreationDate = () => this.#row !== undefined ? this.#row.created : ""
  getUpdateDate = () => this.#row !== undefined ? this.#row.updated : ""
}