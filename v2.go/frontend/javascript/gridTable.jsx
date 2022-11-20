// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class GridTable extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm table-responsive">
				<thead>
					<tr>
						<th scope="col" style={{width: "24px"}}></th>
						{this.props.columns && this.props.columns.map(
							col => <th scope="col" key={col.name}>{col.label}<br/><small>{col.name}</small></th>
						)}
					</tr>
				</thead>
				<tbody className="table-group-divider">
					{this.props.rows.map(
						row => <GridRow key={row.uuid}
										row={row}
										rowSelected={this.props.rowsSelected.includes(row.uuid)}
										rowEdited={this.props.rowsEdited.includes(row.uuid)}
										rowAdded={this.props.rowsAdded.includes(row.uuid)}
										referencedValuesAdded={this.props.referencedValuesAdded.filter(ref => ref.rowUuid == row.uuid)}
										referencedValuesRemoved={this.props.referencedValuesRemoved.filter(ref => ref.rowUuid == row.uuid)}
										columns={this.props.columns}
										onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										onDeleteRowClick={uuid => this.props.onDeleteRowClick(uuid)}
										onAddReferencedValueClick={(rowUuid, col, gridUuid, uuid, displayString, path) => this.props.onAddReferencedValueClick(rowUuid, col, gridUuid, uuid, displayString, path)}
										onRemoveReferencedValueClick={(rowUuid, col, gridUuid, uuid, displayString, path) => this.props.onRemoveReferencedValueClick(rowUuid, col, gridUuid, uuid, displayString, path)}
										inputRef={this.props.inputRef}
										dbName={this.props.dbName}
										token={this.props.token} />
					)}
				</tbody>
			</table>
		)
	}
}

class GridRow extends React.Component {
	render() {
		const variant = this.props.rowEdited ? "table-warning" : ""
		const columns = getColumnValuesForRow(this.props.columns, this.props.row)
		return (
			<tr className={variant}>
				<td scope="row" className="vw-10">
					{!(this.props.rowAdded || this.props.rowSelected) && 
						<a href={this.props.row.path}><i className="bi bi-card-text"></i></a>
					}
					{(this.props.rowAdded || this.props.rowSelected) && 
						<button
							type="button"
							className="btn text-danger btn-sm mx-0 p-0"
							onClick={() => this.props.onDeleteRowClick(this.props.row.uuid)}>
							<i className="bi bi-dash-circle"></i>
						</button>
					}
				</td>
				{columns.map(
					col => <GridCell uuid={this.props.row.uuid}
										key={col.name}
										col={col.name}
										type={col.type}
										typeUuid={col.typeUuid}
										value={col.value}
										values={col.values}
										gridPromptUuid={col.gridPromptUuid}
										gridPromptUri={col.gridPromptUri}
										readonly={col.readonly}
										rowAdded={this.props.rowAdded}
										rowSelected={this.props.rowSelected}
										rowEdited={this.props.rowEdited}
										referencedValuesAdded={this.props.referencedValuesAdded.filter(ref => ref.col == col.name)}
										referencedValuesRemoved={this.props.referencedValuesRemoved.filter(ref => ref.col == col.name)}
										onSelectRowClick={uuid => this.props.onSelectRowClick(uuid)}
										onEditRowClick={uuid => this.props.onEditRowClick(uuid)}
										onAddReferencedValueClick={(rowUuid, col, gridUuid, uuid, displayString, path) => this.props.onAddReferencedValueClick(rowUuid, col, gridUuid, uuid, displayString, path)}
										onRemoveReferencedValueClick={(rowUuid, col, gridUuid, uuid, displayString, path) => this.props.onRemoveReferencedValueClick(rowUuid, col, gridUuid, uuid, displayString, path)}
										inputRef={this.props.inputRef}
										dbName={this.props.dbName}
										token={this.props.token} />
				)}
			</tr>
		)
	}
}

class GridCell extends React.Component {
	render() {
		const variant = this.props.readonly ? " form-control form-control-sm form-control-plaintext rounded-2 shadow " : "form-control form-control-sm rounded-2 shadow "
		const variantSize = this.props.typeUuid == UuidUuidColumnType ? " font-monospace " : ""
		return (
			<td onClick={() => this.props.onSelectRowClick(this.props.uuid)}>
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" && 
					<GridCellInput type={this.props.type}
									variant={variant}
									uuid={this.props.uuid}
									col={this.props.col}
									readOnly={this.props.readonly}
									value={this.props.value}
									inputRef={this.props.inputRef}
									onEditRowClick={uuid => this.props.onEditRowClick(uuid)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type != "reference" &&
					<span className={variantSize}>{getCellValue(this.props.type, this.props.value)}</span>
				}
				{(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" && 
					<GridCellDropDown uuid={this.props.uuid}
										col={this.props.col}
										values={this.props.values}
										dbName={this.props.dbName}
										token={this.props.token}
										gridPromptUuid={this.props.gridPromptUuid}
										gridPromptUri={this.props.gridPromptUri}
										referencedValuesAdded={this.props.referencedValuesAdded}
										referencedValuesRemoved={this.props.referencedValuesRemoved}
										onAddReferencedValueClick={(rowUuid, col, gridUuid, uuid, displayString, path) => this.props.onAddReferencedValueClick(rowUuid, col, gridUuid, uuid, displayString, path)}
										onRemoveReferencedValueClick={(rowUuid, col, gridUuid, uuid, displayString, path) => this.props.onRemoveReferencedValueClick(rowUuid, col, gridUuid, uuid, displayString, path)} />
				}
				{!(this.props.rowAdded || this.props.rowEdited || this.props.rowSelected) && this.props.type == "reference" &&
					<GridCellReferences uuid={this.props.uuid}
										values={this.props.values}
										referencedValuesAdded={this.props.referencedValuesAdded}
										referencedValuesRemoved={this.props.referencedValuesRemoved} />
				}
			</td>
		)
	}
}

class GridCellInput  extends React.Component {
	render() {
		return (
			<input type={this.props.type}
					className={this.props.variant}
					uuid={this.props.uuid}
					col={this.props.col}
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
			filter((value, index, self) => index === self.findIndex((t) => (t.uuid === value.uuid)))
		if(referencedValuesIncluded.length > 0) {
			return (
				<ul className="list-unstyled mb-0">
					{this.props.values.map(value => 
						<li key={value.uuid}>
							<a className="gap-2 p-0" href={value.path}>{value.displayString}</a>
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
		const { isLoading, isLoaded, error, rows } = this.state
		const { values, gridPromptUri, referencedValuesAdded, referencedValuesRemoved } = this.props
		const referencedValuesIncluded = values.
			concat(referencedValuesAdded).
			filter(ref => !referencedValuesRemoved.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index === self.findIndex((t) => (t.uuid === value.uuid)))
		const referencedValuesNotIncluded = rows.
			concat(referencedValuesRemoved).
			filter(ref => !referencedValuesAdded.map(ref => ref.uuid).includes(ref.uuid)).
			filter(ref => !referencedValuesIncluded.map(ref => ref.uuid).includes(ref.uuid)).
			filter((value, index, self) => index === self.findIndex((t) => (t.uuid === value.uuid)))
		const countRows = referencedValuesNotIncluded ? referencedValuesNotIncluded.length : 0
		return (
			<ul className="list-unstyled mb-0">
				{referencedValuesIncluded.map(ref => 
					<li key={ref.uuid}>
						<span>
							<button type="button" className="btn text-danger btn-sm mx-0 p-0"
									onClick={() => this.props.onRemoveReferencedValueClick(this.props.uuid, this.props.col, this.props.gridPromptUuid, ref.uuid, ref.displayString, ref.path)}>
								<i className="bi bi-box-arrow-down pe-1"></i>
							</button>
							{ref.displayString}
						</span>
					</li>
				)}
				<li>
					<input type="search"
							className="form-control form-control-sm rounded-2 shadow gap-2 p-1"
							autoComplete="false"
							placeholder="Search..."
							onInput={(e) => this.loadDropDownData(gridPromptUri, e.target.value)} />
				</li>
				{isLoading && <li><Spinner /></li>}
				{error && !isLoading && !isLoaded && <li className="alert alert-danger" role="alert">{error}</li>}
				{error && !isLoading && isLoaded && <li className="alert alert-primary" role="alert">{error}</li>}
				{referencedValuesNotIncluded && countRows > 0 && referencedValuesNotIncluded.map(ref => (
					<li key={ref.uuid}>
						<span>
							<button type="button"
									className="btn text-success btn-sm mx-0 p-0"
									onClick={() => this.props.onAddReferencedValueClick(this.props.uuid, this.props.col, this.props.gridPromptUuid, ref.uuid, ref.displayString, ref.path)}>
								<i className="bi bi-box-arrow-up pe-1"></i>
							</button>
							{ref.displayString}
						</span>
					</li>
				))}
			</ul>
		)
	}

	loadDropDownData(gridPromptUri, value) {
		this.setState({isLoading: true})
		if(value.length > 0) {
			const uri = `/${this.props.dbName}/api/v1/${gridPromptUri}?trace=true`
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
							this.setState({
								isLoading: false,
								isLoaded: true,
								rows: result.rows,
								error: result.error
							})
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
				rows: []
			})
		}
	}
}
