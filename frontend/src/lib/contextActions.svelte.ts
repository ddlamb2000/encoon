// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2025

import type { GridResponse, RowType, ColumnType, GridType, ReferenceType } from '$lib/apiTypes'
import { ContextData } from '$lib/contextData.svelte.ts'
import { newUuid, debounce, numberToLetters } from "$lib/utils.svelte.ts"
import { replaceState } from "$app/navigation"
import * as metadata from "$lib/metadata.svelte"

export class ContextActions extends ContextData {
  load = async () => {
    this.pushTransaction({
      action: metadata.ActionLoad,
      actionText: "Load " + (this.uuid !== "" ? "row" : "grid"),
      gridUuid: this.gridUuid,
      uuid: this.uuid
    })
  }

  async changeFocus(grid: GridType | undefined, column: ColumnType | undefined, row: RowType | undefined) {
    console.log("changeFocus[1]", row !== undefined ? row.uuid : undefined)
    if(grid) {
      console.log("changeFocus[2]", row !== undefined ? row.uuid : undefined)
      await this.pushTransaction(
        {
          action: metadata.ActionLocateGrid,
          gridUuid: grid.uuid,
          columnUuid: column !== undefined ? column.uuid : undefined,
          uuid: row !== undefined ? row.uuid : undefined
        }
      )
    }
  }

  navigateToGrid = async (gridUuid: string, uuid?: string) => {
		console.log(`[Context.navigateToGrid()] gridUuid=${gridUuid}, uuid=${uuid}`)
    this.reset()
    const url = `/${this.dbName}/${gridUuid}` + (uuid !== "" ? `/${uuid}` : "")
    replaceState(url, { gridUuid: this.gridUuid, uuid: this.uuid })
    this.gridUuid = gridUuid
    this.uuid = uuid ?? ""
    this.load()
	}

  changeCell = debounce(
    async (set: GridResponse, row: RowType) => {
      row.updated = new Date
      const rowClone = Object.assign({}, row)
      if(set.grid.columns) {
        for(const column of set.grid.columns) {
          if(column.typeUuid === metadata.UuidIntColumnType) {
            if(!row[column.name] || row[column.name] === "" || row[column.name] === "<br>") rowClone[column.name] = undefined
            else if(typeof row[column.name] === "string") rowClone[column.name] = row[column.name].replace(/[^0-9-]/g, "") * 1
          }
        }
      }
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: 'Update',
          gridUuid: set.grid.uuid,
          dataSet: { rowsEdited: [rowClone] }
        }
      )
    },
    500
  )

  getPrefixFromColumknType = (columnTypeUuid: string): string => {
    switch(columnTypeUuid) {
      case metadata.UuidTextColumnType:
      case metadata.UuidPasswordColumnType:
      case metadata.UuidUuidColumnType:
      case metadata.UuidBooleanColumnType:
        return "text"
      case metadata.UuidIntColumnType:
        return "int"
      case metadata.UuidReferenceColumnType:
        return "relationship"
      }
    return ""
  }

  getColumnName = (set: GridResponse, rowPrompt: RowType): string => {
    if(set.grid && rowPrompt.uuid) {
      const prefixColumnName = this.getPrefixFromColumknType(rowPrompt.uuid)
      const columnsSamePrefix = set.grid.columns !== undefined ? set.grid.columns.filter((c) => this.getPrefixFromColumknType(c.typeUuid) === prefixColumnName) : undefined
      const nbColumnsSamePrefix = columnsSamePrefix !== undefined ? columnsSamePrefix.length : 0
      if(nbColumnsSamePrefix < 10) {
        const columnName = prefixColumnName + (nbColumnsSamePrefix + 1)
        return columnName
      }
    }
    return ""
  }

  addColumn = async (set: GridResponse, rowPrompt: RowType, rowReference: RowType | undefined = undefined) => {
    const uuidColumn = newUuid()
    const nbColumns = set.grid.columns ? set.grid.columns.length : 0
    const newLabel = numberToLetters(nbColumns)
    const columnName = this.getColumnName(set, rowPrompt)
    if(columnName !== "") {
      const column: ColumnType = { uuid: uuidColumn,
                                    orderNumber: nbColumns + 1,
                                    owned: true,
                                    label: newLabel,
                                    name: columnName,
                                    type: rowPrompt.text1 || "?",
                                    typeUuid: rowPrompt.uuid,
                                    gridUuid: set.grid.uuid,
                                    gridPromptUuid: rowReference !== undefined ? rowReference.uuid : undefined
                                  }
      if(set.grid.columns) set.grid.columns.push(column)
      else set.grid.columns = [column]
      const rowsAdded = [
        { gridUuid: metadata.UuidColumns,
          uuid: uuidColumn,
          text1: newLabel,
          text2: columnName,
          int1: nbColumns + 1,
          created: new Date,
          updated: new Date } 
      ]
      const referencedValuesAdded = [
        { owned: false,
          columnName: "relationship1",
          fromUuid: uuidColumn,
          toGridUuid: metadata.UuidGrids,
          uuid: set.grid.uuid },
        { owned: true,
          columnName: "relationship1",
          fromUuid: uuidColumn,
          toGridUuid: metadata.UuidColumnTypes,
          uuid: rowPrompt.uuid }
      ] 
      if(rowReference !== undefined) {
        referencedValuesAdded.push(
          { owned: true,
            columnName: "relationship2",
            fromUuid: uuidColumn,
            toGridUuid: metadata.UuidGrids,
            uuid: rowReference.uuid }  
        )
      }
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: 'Add column',
        gridUuid: metadata.UuidColumns,
        dataSet: { rowsAdded: rowsAdded, referencedValuesAdded: referencedValuesAdded }
      })
    }
  }
  
  addRow = async (set: GridResponse) => {
    const uuid = newUuid()
    const row: RowType = { gridUuid: set.grid.uuid, uuid: uuid, created: new Date, updated: new Date }
    set.rows.push(row)
    set.countRows += 1
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'Add row',
      gridUuid: set.grid.uuid,
      dataSet: { rowsAdded: [row] }
    })
  }

  removeRow = async (set: GridResponse, row: RowType) => {
    const rowIndex = set.rows.findIndex((r) => r.uuid === row.uuid)
    if(rowIndex >= 0) {
      const deletedRow: RowType = { gridUuid: set.grid.uuid, uuid: row.uuid }
      set.rows.splice(rowIndex, 1)
      set.countRows -= 1
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: 'Remove row',
        gridUuid: set.grid.uuid,
        dataSet: { rowsDeleted: [deletedRow] }
      })
    }
  }

  removeColumn = async (set: GridResponse, column: ColumnType) => {
    if(set.grid.columns && set.grid.columns !== undefined && column !== undefined && column.uuid !== undefined) {
      const columnIndex = set.grid.columns.findIndex((c) => c.uuid === column.uuid)
      set.grid.columns.splice(columnIndex, 1)
      return this.pushTransaction({
        action: metadata.ActionChangeGrid,
        actionText: 'Remove column',
        gridUuid: metadata.UuidColumns,
        dataSet: {
          rowsDeleted: [
            { gridUuid: metadata.UuidColumnTypes,
              uuid: column.uuid }
          ],
          referencedValuesRemoved: [
            { owned: false,
              columnName: "relationship1",
              fromUuid: column.uuid,
              toGridUuid: metadata.UuidGrids,
              uuid: set.grid.uuid },
            { owned: true,
              columnName: "relationship1",
              fromUuid: column.uuid,
              toGridUuid: metadata.UuidColumnTypes,
              uuid: column.typeUuid }
          ] 
        }
      })
    }
  }

  newGrid = async () => {
    const gridUuid = newUuid()
    await this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'New grid',
      gridUuid: metadata.UuidGrids,
      dataSet: {
        rowsAdded: [
          { gridUuid: metadata.UuidGrids,
            uuid: gridUuid,
            displayString: 'New grid',
            text1: 'New grid',
            text2: 'Untitled',
            text3: 'journal',
            created: new Date,
            updated: new Date } 
        ]
      }
    })

    const uuidColumn = newUuid()
    await this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'Add column into new grid',
      gridUuid: metadata.UuidColumns,
      dataSet: { rowsAdded: [
        { gridUuid: metadata.UuidColumns,
          uuid: uuidColumn,
          text1: 'A',
          text2: 'text1',
          int1: 1,
          created: new Date,
          updated: new Date } 
      ], referencedValuesAdded: [
        { owned: false,
          columnName: "relationship1",
          fromUuid: uuidColumn,
          toGridUuid: metadata.UuidGrids,
          uuid: gridUuid },
        { owned: true,
          columnName: "relationship1",
          fromUuid: uuidColumn,
          toGridUuid: metadata.UuidColumnTypes,
          uuid: metadata.UuidTextColumnType }
      ] 
     }
    })

    const uuid = newUuid()
    const row: RowType = { gridUuid: gridUuid, uuid: uuid, created: new Date, updated: new Date }
    await this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'Add row into new grid',
      gridUuid: gridUuid,
      dataSet: { rowsAdded: [
        { gridUuid: gridUuid, uuid: uuid, created: new Date, updated: new Date }
      ] }
    })

    this.navigateToGrid(gridUuid, "")
  }

  addReferencedValue = async (set: GridResponse, column: ColumnType, row: RowType, rowPrompt: RowType) => {
    const reference = row.references !== undefined ? 
                        row.references.find((reference) => reference.owned && reference.name === column.name) :
                        undefined
    if(reference !== undefined) {
      if(reference.rows !== undefined) reference.rows.push(rowPrompt)
      else reference.rows = [rowPrompt]
    } else {
      const reference: ReferenceType = {
        owned: true,
        label: column.label,
        name: column.name,
        gridUuid: column.gridPromptUuid,
        rows: [rowPrompt]
      }
      if(row.references !== undefined) row.references.push(reference)
      else row.references = [reference]
    }
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'Add value',
      gridUuid: set.grid.uuid,
      dataSet: {
        referencedValuesAdded: [
          { owned: true,
            columnName: column.name,
            fromUuid: row.uuid,
            toGridUuid: rowPrompt.gridUuid,
            uuid: rowPrompt.uuid },
        ] 
      }
    })    
  }

  removeReferencedValue = async (set: GridResponse, column: ColumnType, row: RowType, rowPrompt: RowType) => {
    if(row.references !== undefined) {
      const reference = row.references.find((reference) => reference.owned && reference.name === column.name)
      if(reference !== undefined) {
        if(reference.rows !== undefined) {
          const rowIndex = reference.rows.findIndex((r) => r.uuid === rowPrompt.uuid)
          if(rowIndex >= 0) reference.rows.splice(rowIndex, 1)
        }
      }
    }
    return this.pushTransaction({
      action: metadata.ActionChangeGrid,
      actionText: 'Remove value',
      gridUuid: set.grid.uuid,
      dataSet: {
        referencedValuesRemoved: [
          { owned: true,
            columnName: column.name,
            fromUuid: row.uuid,
            toGridUuid: rowPrompt.gridUuid,
            uuid: rowPrompt.uuid },
        ] 
      }
    })    
  }

  changeGrid = debounce(
    async (grid: GridType) => {
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: 'Update grid',
          gridUuid: metadata.UuidGrids,
          dataSet: { rowsEdited: [grid] }
        }
      )
    },
    500
  )

  changeColumn = debounce(
    async (grid: GridType, column: ColumnType) => {
      this.pushTransaction(
        {
          action: metadata.ActionChangeGrid,
          actionText: 'Update column',
          gridUuid: metadata.UuidColumns,
          dataSet: {
            rowsEdited: [
              { gridUuid: metadata.UuidColumns,
                uuid: column.uuid,
                text1: column.label,
                text2: column.name,
                int1: column.orderNumber,
                updated: new Date }               
            ] 
          }
        }
      )
    },
    500
  )

  locateGrid = (gridUuid: string | undefined, columnUuid: string | undefined, uuid: string | undefined) => {
    console.log(`[Context.locateGrid(${gridUuid},${columnUuid},${uuid})`)
    if(gridUuid) {
      for(const set of this.dataSet) {
        if(set && set.grid && set.gridUuid === gridUuid) {
          const grid: GridType = set.grid
          if(grid.columns) {
            const column: ColumnType | undefined = grid.columns.find((column) => column.uuid === columnUuid)
            if(column) {
              const row = set.rows.find((row) => row.uuid === uuid)
              this.focus.set(grid, column, row)
              return
            }
            else {
              const row = set.rows.find((row) => row.uuid === uuid)
              this.focus.set(grid, undefined, row)
              return
            }
          } else {
            this.focus.set(grid, undefined, undefined)
            return
          }
        }
      }
    }
    this.focus.reset()
  }
}