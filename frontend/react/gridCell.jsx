// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class GridCell extends React.Component {
	render() {
		const variantReadOnly = this.props.readonly ? "form-control-plaintext" : ""
		const checkedBoolean = this.props.value && this.props.value == "true" ? true : false
		const variantEdited = this.props.rowEdited ? "table-warning" : ""
		const embedded = this.props.typeUuid == UuidReferenceColumnType && this.props.bidirectional && this.props.owned
		const id = this.props.columnName + '-' + this.props.uuid + '-' + this.props.columnUuid
		return (
			<td className={variantEdited} onClick={() => this.props.onSelectRowClick(this.props.canEditRow ? this.props.uuid : '')}>
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.typeUuid != UuidReferenceColumnType && !embedded && 
					<GridCellInput id={id}
									type={this.props.type}
									typeUuid={this.props.typeUuid}
									variantReadOnly={variantReadOnly}
									uuid={this.props.uuid}
									columnUuid={this.props.columnUuid}
									columnName={this.props.columnName}
									readOnly={this.props.readonly}
									checkedBoolean={checkedBoolean}
									value={this.props.value}
									inputRef={this.props.inputRef}
									onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
									createRichTextField={(id, value, display) => this.props.createRichTextField(id, value, display)}
									deleteRichTextField={id => this.props.deleteRichTextField(id)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.typeUuid != UuidReferenceColumnType && !embedded &&
					<GridCellDisplay id={id}
										typeUuid={this.props.typeUuid}
										value={this.props.value}
										icon={this.props.icon}
										createRichTextField={(id, value, display) => this.props.createRichTextField(id, value, display)}
										deleteRichTextField={id => this.props.deleteRichTextField(id)} />
				}
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.typeUuid == UuidReferenceColumnType && !embedded && 
					<GridCellDropDown uuid={this.props.uuid}
										columnUuid={this.props.columnUuid}
										columnName={this.props.columnName}
										owned={this.props.owned}
										values={this.props.values}
										dbName={this.props.dbName}
										token={this.props.token}
										gridPromptUuid={this.props.gridPromptUuid}
										referencedValuesAdded={this.props.referencedValuesAdded}
										referencedValuesRemoved={this.props.referencedValuesRemoved}
										onAddReferencedValueClick={reference => this.props.onAddReferencedValueClick(reference)}
										onRemoveReferencedValueClick={reference => this.props.onRemoveReferencedValueClick(reference)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.typeUuid == UuidReferenceColumnType && !embedded &&
					<GridCellReferences uuid={this.props.uuid}
										values={this.props.values}
										referencedValuesAdded={this.props.referencedValuesAdded}
										referencedValuesRemoved={this.props.referencedValuesRemoved}
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
				}
				{embedded && 
					<Grid token={this.props.token}
							dbName={this.props.dbName} 
							gridUuid={this.props.gridPromptUuid}
							filterColumnOwned={this.props.owned ? 'false' : 'true'}
							filterColumnName={this.props.columnName}
							filterColumnLabel={this.props.columnLabel}
							filterColumnGridUuid={this.props.grid.uuid}
							filterColumnValue={this.props.uuid}
							filterColumnDisplayString={this.props.displayString}
							navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
							innerGrid={true}
							loadParentData={() => this.props.loadParentData()} />
				}
			</td>
		)
	}
}
