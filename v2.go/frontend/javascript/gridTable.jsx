// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm table-responsive align-middle">
				<thead>
					<tr>
						<th scope="col" style={{width: "24px"}}></th>
						{this.props.columns && this.props.columns.map(
							column =>
								<th scope="col" key={column.uuid}>
									{column.label}
									{!column.owned && <small><br />{column.grid.displayString}</small>}
								</th>
						)}
						<th className="text-end" scope="col">Revision</th>
					</tr>
				</thead>
				<tbody className="table-group-divider">
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
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
					)}
				</tbody>
			</table>
		)
	}
}

class GridRow extends React.Component {
	render() {
		const { row, rowAdded, rowSelected, rowEdited, referencedValuesAdded, referencedValuesRemoved } = this.props
		const variantEdited = this.props.rowEdited ? "table-warning" : ""
		const columns = getColumnValuesForRow(this.props.columns, row)
		return (
			<tr className={variantEdited}>
				<td scope="row" className="vw-10">
					{!(rowAdded || rowSelected) && 
						<a href="#" onClick={() => this.props.navigateToGrid(row.gridUuid, row.uuid)}>
						<i className="bi bi-card-text"></i>
					</a>
					}
					{row.canEditRow && (rowAdded || rowSelected) && 
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
										columnName={column.name}
										type={column.type}
										typeUuid={column.typeUuid}
										value={column.value}
										values={column.values}
										gridPromptUuid={column.gridPromptUuid}
										readonly={column.readonly}
										canEditRow={row.canEditRow}
										rowAdded={rowAdded}
										rowSelected={rowSelected}
										rowEdited={rowEdited}
										referencedValuesAdded={referencedValuesAdded.filter(ref => ref.columnUuid == column.uuid)}
										referencedValuesRemoved={referencedValuesRemoved.filter(ref => ref.columnUuid == column.uuid)}
										onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										onAddReferencedValueClick={reference => this.props.onAddReferencedValueClick(reference)}
										onRemoveReferencedValueClick={reference => this.props.onRemoveReferencedValueClick(reference)}
										inputRef={this.props.inputRef}
										dbName={this.props.dbName}
										token={this.props.token}
										grid={this.props.grid}
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
				)}
				<td className="text-end">{this.props.row.revision}</td>
			</tr>
		)
	}
}

class GridCell extends React.Component {
	render() {
		const variantReadOnly = this.props.readonly ? "form-control-plaintext" : ""
		const variantMonospace = this.props.typeUuid == UuidUuidColumnType ? " font-monospace " : ""
		return (
			<td onClick={() => this.props.onSelectRowClick(this.props.canEditRow ? this.props.uuid : '')}>
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" && 
					<GridCellInput type={this.props.type}
									variantReadOnly={variantReadOnly}
									uuid={this.props.uuid}
									columnUuid={this.props.columnUuid}
									columnName={this.props.columnName}
									readOnly={this.props.readonly}
									value={this.props.value}
									inputRef={this.props.inputRef}
									onEditRowClick={uuid => this.props.onEditRowClick(uuid)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" &&
					<span className={variantMonospace}>{getCellValue(this.props.type, this.props.value)}</span>
				}
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" && 
					<GridCellDropDown uuid={this.props.uuid}
										columnUuid={this.props.columnUuid}
										columnName={this.props.columnName}
										values={this.props.values}
										dbName={this.props.dbName}
										token={this.props.token}
										gridPromptUuid={this.props.gridPromptUuid}
										referencedValuesAdded={this.props.referencedValuesAdded}
										referencedValuesRemoved={this.props.referencedValuesRemoved}
										onAddReferencedValueClick={reference => this.props.onAddReferencedValueClick(reference)}
										onRemoveReferencedValueClick={reference => this.props.onRemoveReferencedValueClick(reference)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" &&
					<GridCellReferences uuid={this.props.uuid}
										values={this.props.values}
										referencedValuesAdded={this.props.referencedValuesAdded}
										referencedValuesRemoved={this.props.referencedValuesRemoved}
										navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)} />
				}
			</td>
		)
	}
}

class GridCellInput  extends React.Component {
	render() {
		return (
			<input type={this.props.type}
					className={"form-control form-control-sm rounded-2 shadow " + this.props.variantReadOnly}
					uuid={this.props.uuid}
					column={this.props.columnName}
					readOnly={this.props.readonly}
					defaultValue={this.props.value}
					ref={this.props.inputRef}
					onInput={() => this.props.onEditRowClick(this.props.uuid)} />
		)
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
				if(contentType && contentType.indexOf("application/json") !== -1) {
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
