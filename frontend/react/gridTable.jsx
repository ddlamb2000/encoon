// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm table-responsive align-middle">
				{!this.props.miniGrid && 
					<thead>
						<tr>
							{<th scope="col" style={{width: "24px"}}></th>}
							{this.props.columns && this.props.columns.map( 
								column => <GridRowHeader key={column.uuid}
															column={column}
															filterColumnOwned={this.props.filterColumnOwned}
															filterColumnName={this.props.filterColumnName}
															filterColumnGridUuid={this.props.filterColumnGridUuid}
															grid={this.props.grid} />
							)}
						</tr>
					</thead>
				}
				<tbody className={!this.props.miniGrid ? "table-group-divider" : ""}>
					{this.props.rows.map(
						row => <GridRow key={row.uuid}
										row={row}
										rowSelected={this.props.rowsSelected.includes(row.uuid)}
										rowEdited={this.props.rowsEdited.includes(row.uuid)}
										rowAdded={this.props.rowsAdded.includes(row.uuid)}
										referencedValuesAdded={this.props.referencedValuesAdded.filter(ref => ref.fromUuid == row.uuid)}
										referencedValuesRemoved={this.props.referencedValuesRemoved.filter(ref => ref.fromUuid == row.uuid)}
										columns={this.props.columns}
										onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										onDeleteRowClick={uuid => this.props.onDeleteRowClick(uuid)}
										onAddReferencedValueClick={reference => this.props.onAddReferencedValueClick(reference)}
										onRemoveReferencedValueClick={reference => this.props.onRemoveReferencedValueClick(reference)}
										inputRef={this.props.inputRef}
										dbName={this.props.dbName}
										token={this.props.token}
										grid={this.props.grid}
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
										createRichTextField={(id, value, display) => this.props.createRichTextField(id, value, display)}
										deleteRichTextField={id => this.props.deleteRichTextField(id)}
										miniGrid={this.props.miniGrid} />
					)}
				</tbody>
			</table>
		)
	}
}
