// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

function getHtmlColumnType(type) {
	switch(type) {
		case UuidIntColumnType:
			return "number"
		case UuidPasswordColumnType:
			return "password"
		case UuidBooleanColumnType:
			return "checkbox"
		default:
			return "text"
	}
}

function getColumnValuesForRow(columns, row, withTimeStamps) {
	const cols = []
	{columns && columns.map(
		column => {
			let type = getHtmlColumnType(column.typeUuid)
			if(column.typeUuid == UuidReferenceColumnType) {
				cols.push({
					uuid: column.uuid,
					owned: column.owned,
					name: column.name,
					label: column.label,
					bidirectional: column.bidirectional,
					typeUuid: column.typeUuid,
					gridPromptUuid: column.gridPromptUuid,
					gridUuid: column.gridUuid,
					grid: column.grid,
					typeUuid: column.typeUuid,
					type: type,
					readonly: false,
					values: getColumnValueForReferencedRow(column, row)
				})
			} else {
				cols.push({
					uuid: column.uuid,
					owned: column.owned,
					name: column.name,
					label: column.label,
					value: row[column.name],
					typeUuid: column.typeUuid,
					gridUuid: column.gridUuid,
					grid: column.grid,
					typeUuid: column.typeUuid,
					type: type,
					readonly: false
				})	
			}
		}
	)}
	if(withTimeStamps) {
		cols.push({uuid: "a", name: "uuid", label: "Identifier", value: row.uuid, typeUuid: UuidUuidColumnType, typeUuid: UuidUuidColumnType, type: "text", owned: true, readonly: true})
		cols.push({uuid: "b", name: "revision", label: "Revision", value: row.revision, typeUuid: UuidIntColumnType, type: "number", owned: true, readonly: true})
	}
	return cols
}

function getColumnValueForReferencedRow(column, row) {
	let output = []
	if(row.references) {
		row.references.map(
			ref => {
				if(ref.gridUuid == column.gridUuid && ref.name == column.name && ref.rows) {
					ref.rows.map(
						refRow => output.push({
							gridUuid: refRow.gridUuid,
							uuid: refRow.uuid,
							displayString: refRow.displayString
						})
					)
				}
			}
		)
	}
	return output
}

function matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) {
    const ownership = (column.owned && filterColumnOwned == 'true') || (!column.owned && filterColumnOwned != 'true')
    return ownership && column.name == filterColumnName && column.gridUuid == filterColumnGridUuid
}
