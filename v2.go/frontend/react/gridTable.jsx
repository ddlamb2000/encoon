// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

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

class GridRowHeader extends React.Component {
	render() {
		const { column, filterColumnOwned, filterColumnName, filterColumnGridUuid } = this.props
		return (
			<th scope="col">
				{this.matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) && <mark>{column.label}</mark>}
				{!this.matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) && column.label}
				{!column.owned && <small><br />[{column.grid.displayString}]</small>}
				{trace && <small><br /><em>{column.gridUuid}</em></small>}
				{trace && <small><br /><em>{column.name}</em></small>}
			</th>
		)
	}

	matchFilter(column, filterColumnOwned, filterColumnName, filterColumnGridUuid) {
		const ownership = (column.owned && filterColumnOwned == 'true') || (!column.owned && filterColumnOwned != 'true')
		return ownership && column.name == filterColumnName && column.gridUuid == filterColumnGridUuid
	}
}

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

class GridCellDisplay extends React.Component {
	render() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			return (
				<div id={"richtext-" + this.props.id}>
					<div id={this.props.id}
							typeuuid={this.props.typeUuid}
							uuid={this.props.uuid}
							ref={this.props.inputRef}
							column={this.props.columnName}
							onInput={() => this.props.onEditRowClick(this.props.uuid)} />
				</div>
			)
		}
		else {
			const variantMonospace = this.props.typeUuid == UuidUuidColumnType ? " font-monospace " : ""
			return (
				<span className={variantMonospace}>
					{this.getCellValue(this.props.typeUuid, this.props.value)} {this.props.icon && <i className={`bi bi-${this.props.icon}`}></i>}
				</span>
			)
		}
	}

	getCellValue(typeUuid, value) {
		switch(typeUuid) {
			case UuidPasswordColumnType: return '*****'
			case UuidBooleanColumnType: return value == 'true' ? '✔︎' : ''
		}
		return value
	}

	componentDidMount() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.createRichTextField(this.props.id, this.props.value, true)
		}
	}

	componentWillUnmount() {
		if(trace) console.log("[GridCellDisplay.componentWillUnmount()]")
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.deleteRichTextField(this.props.id)
		}
	}
}

class GridCellInput extends React.Component {
	render() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			return (
				<div id={"richtext-" + this.props.id}>
					<div id={this.props.id}
							typeuuid={this.props.typeUuid}
							uuid={this.props.uuid}
							ref={this.props.inputRef}
							column={this.props.columnName}
							onInput={() => this.props.onEditRowClick(this.props.uuid)} />
				</div>
			)
		}
		else {
			const className = this.props.type == "checkbox" ? "form-check-input" : "form-control"
			return (
				<input type={this.props.type}
						typeuuid={this.props.typeUuid}
						className={className + " form-control-sm rounded-2 shadow gap-2 p-1 " + this.props.variantReadOnly}
						name={this.props.uuid}
						uuid={this.props.uuid}
						column={this.props.columnName}
						readOnly={this.props.readonly}
						defaultChecked={this.props.checkedBoolean}
						defaultValue={this.props.value}
						ref={this.props.inputRef}
						onInput={() => this.props.onEditRowClick(this.props.uuid)} />
			)
		}
	}

	componentDidMount() {
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.createRichTextField(this.props.id, this.props.value, false)
		}
	}

	componentWillUnmount() {
		if(trace) console.log("[GridCellInput.componentWillUnmount()]")
		if(this.props.typeUuid == UuidRichTextColumnType) {
			this.props.deleteRichTextField(this.props.id)
		}
	}
}

class GridCellReferences extends React.Component {
	render() {
		const { values, referencedValuesAdded, referencedValuesRemoved } = this.props
		const referencedValuesIncluded = values.
			concat(referencedValuesAdded).
			filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index == self.findIndex((t) => (t.uuid == value.uuid)))
		if(referencedValuesIncluded.length > 0) {
			return (
				<ul className="list-unstyled mb-0">
					{values.map(value => 
						<li key={value.uuid}>
							<a className="gap-2 p-0" href="#" onClick={() => this.props.navigateToGrid(value.gridUuid, value.uuid)}>
								{value.displayString}
							</a>
						</li>
					)}
				</ul>
			)
		}
	}
}

class GridCellDropDown extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: "",
			isLoaded: false,
			isLoading: false,
			rows: [],
		}
	}

	render() {
		const { isLoading, error, rows } = this.state
		const { gridPromptUuid, values, referencedValuesAdded, referencedValuesRemoved } = this.props
		const referencedValuesIncluded = error ? [] : values.
			concat(referencedValuesAdded).
			filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index == self.findIndex((t) => (t.uuid == value.uuid)))
		const referencedValuesNotIncluded = error ? [] : rows.
			concat(referencedValuesRemoved).
			filter(ref => !referencedValuesAdded.map(ref => ref.uuid).includes(ref.uuid)).
			filter(ref => !referencedValuesIncluded.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index == self.findIndex((t) => (t.uuid == value.uuid)))
		const countRows = referencedValuesNotIncluded ? referencedValuesNotIncluded.length : 0
		if(trace) {
			console.log("[GridCellDropDown.render()] this.props.columnName=", this.props.columnName)
			console.log("[GridCellDropDown.render()] values=", values)
			console.log("[GridCellDropDown.render()] referencedValuesAdded=", referencedValuesAdded)
			console.log("[GridCellDropDown.render()] referencedValuesRemoved=", referencedValuesRemoved)
			console.log("[GridCellDropDown.render()] referencedValuesIncluded=", referencedValuesIncluded)
			console.log("[GridCellDropDown.render()] referencedValuesNotIncluded=", referencedValuesNotIncluded)
		}
		return (
			<ul className="list-unstyled mb-0">
				{referencedValuesIncluded.map(ref => 
					<li key={ref.uuid}>
						<span>
							<button type="button" className="btn text-danger btn-sm mx-0 p-0"
									onClick={() => 
										this.props.onRemoveReferencedValueClick({fromUuid: this.props.uuid, 
																				 columnUuid: this.props.columnUuid,
																				 owned: this.props.owned,
																				 columnName: this.props.columnName,
																				 toGridUuid: this.props.gridPromptUuid,
																				 uuid: ref.uuid,
																				 displayString: ref.displayString})}>
								<i className="bi bi-box-arrow-down pe-1"></i>
							</button>
							{ref.displayString}
						</span>
					</li>
				)}
				{gridPromptUuid &&
					<li>
						<input type="search"
								className="form-control form-control-sm rounded-2 shadow gap-2 p-1"
								autoComplete="false"
								placeholder="Search..."
								onInput={(e) => this.loadDropDownData(gridPromptUuid, e.target.value)} />
					</li>
				}
				{isLoading && <li><Spinner /></li>}
				{error && !isLoading && <li className="alert alert-danger" role="alert">{error}</li>}
				{referencedValuesNotIncluded && countRows > 0 && referencedValuesNotIncluded.map(ref => (
					<li key={ref.uuid}>
						<span>
							<button type="button"
									className="btn text-success btn-sm mx-0 p-0"
									onClick={() => this.props.onAddReferencedValueClick({fromUuid: this.props.uuid, 
																						 columnUuid: this.props.columnUuid,
																						 owned: this.props.owned,
																						 columnName: this.props.columnName,
																						 toGridUuid: this.props.gridPromptUuid, 
																						 uuid: ref.uuid, 
																						 displayString: ref.displayString})}>
								<i className="bi bi-box-arrow-up pe-1"></i>
							</button>
							{ref.displayString}
						</span>
					</li>
				))}
			</ul>
		)
	}

	loadDropDownData(gridPromptUuid, value) {
		this.setState({isLoading: true})
		if(value.length > 0) {
			const uri = `/${this.props.dbName}/api/v1/${gridPromptUuid}`
			fetch(uri, {
				headers: {
					'Accept': 'application/json',
					'Content-Type': 'application/json',
					'Authorization': 'Bearer ' + this.props.token
				}
			})
			.then(response => {
				const contentType = response.headers.get("content-type")
				if(contentType && contentType.indexOf("application/json") != -1) {
					return response.json().then(	
						(result) => {
							if(result.response != undefined) {
								this.setState({
									isLoading: false,
									isLoaded: true,
									rows: result.response.rows,
									error: result.response.error
								})
							} else {
								this.setState({
									isLoading: false,
									isLoaded: true,
									error: result.error
								})
							}
						},
						(error) => {
							this.setState({
								isLoading: false,
								isLoaded: false,
								rows: [],
								error: error.message
							})
						}
					)
				} else {
					this.setState({
						isLoading: false,
						isLoaded: false,
						rows: [],
						error: `[${response.status}] Internal server issue.`
					})
				}
			})
		} else {
			this.setState({
				isLoading: false,
				isLoaded: false,
				rows: [],
				error: ""
			})
		}
	}
}
