// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridRow extends React.Component {
	render() {
		const { row, rowAdded, rowSelected, rowEdited, referencedValuesAdded, referencedValuesRemoved, miniGrid } = this.props
		const columns = getColumnValuesForRow(this.props.columns, row, false)
		const icon = row && this.props.grid.uuid == UuidGrids ? row.text3 : ''
		return (
			<tr>
				<td scope="row" className="vw-10">
					{!(rowAdded || rowSelected) && 
						<a href="#" onClick={() => this.props.navigateToGrid(row.gridUuid, row.uuid)}>
							<i className="bi bi-box-arrow-up-right"></i>
						</a>
					}
					{row.canEditRow && (rowAdded || rowSelected) && !miniGrid &&
						<button
							type="button"
							className="btn text-danger btn-sm mx-0 p-0"
							onClick={() => this.props.onDeleteRowClick(row.uuid)}>
							<i className="bi bi-dash-circle"></i>
						</button>
					}
				</td>
				{columns.map(
					column => <GridCell uuid={row.uuid}
										key={column.uuid}
										columnUuid={column.uuid}
										owned={column.owned}
										columnName={column.name}
										columnLabel={column.label}
										type={column.type}
										typeUuid={column.typeUuid}
										value={column.value}
										values={column.values}
										gridPromptUuid={column.gridPromptUuid}
										readonly={column.readonly}
										bidirectional={false}
										canEditRow={row.canEditRow}
										rowAdded={rowAdded}
										rowSelected={rowSelected}
										rowEdited={rowEdited}
										referencedValuesAdded={referencedValuesAdded.filter(ref => ref.columnUuid == column.uuid)}
										referencedValuesRemoved={referencedValuesRemoved.filter(ref => ref.columnUuid == column.uuid)}
										onSelectRowClick={!miniGrid ? uuid => this.props.onSelectRowClick(uuid) : undefined}
										onEditRowClick={!miniGrid ? uuid => this.props.onEditRowClick(uuid) : undefined}
										onAddReferencedValueClick={reference => this.props.onAddReferencedValueClick(reference)}
										onRemoveReferencedValueClick={reference => this.props.onRemoveReferencedValueClick(reference)}
										inputRef={this.props.inputRef}
										dbName={this.props.dbName}
										token={this.props.token}
										grid={this.props.grid}
										icon={column.uuid == UuidGridColumnName ? icon : ''}
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
										createRichTextField={(id, value, display) => this.props.createRichTextField(id, value, display)}
										deleteRichTextField={id => this.props.deleteRichTextField(id)} />
				)}
			</tr>
		)
	}
}
