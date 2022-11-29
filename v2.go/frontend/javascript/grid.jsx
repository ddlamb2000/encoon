// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Grid extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: "",
			isLoaded: false,
			isLoading: false,
			grid: [],
			canAddRows: false,
			rows: [],
			rowsSelected: [],
			rowsEdited: [],
			rowsAdded: [],
			rowsDeleted: [],
			referencedValuesAdded: [],
			referencedValuesRemoved: []
		}
		this.gridInput = new Map()
		this.setGridRowRef = element => {
			if(element != undefined) {
				const uuid = element.getAttribute("uuid")
				const col = element.getAttribute("col")
				if(uuid != "" && col != "") {
					const gridInputMap = this.gridInput.get(uuid)
					if(gridInputMap) {
						gridInputMap.set(col, element)	
					}
					else {
						const gridRowInputMap = new Map()
						gridRowInputMap.set(col, element)
						this.gridInput.set(uuid, gridRowInputMap)
					}
				}
			}
		}
	}

	componentDidMount() {
		this.loadData()
	}

	render() {
		const { isLoading, isLoaded, error, grid, canAddRows, rows, rowsSelected, rowsEdited, rowsAdded, rowsDeleted, referencedValuesAdded, referencedValuesRemoved } = this.state
		const { uuid, dbName, token } = this.props
		const countRows = rows ? rows.length : 0
		return (
			<div className="card my-4">
				<div className="card-body">
					{isLoaded && rows && grid  && uuid == "" && 
						<h4 className="card-title">
							{grid.text1} {grid.text3 && <small><i className={`bi bi-${grid.text3} mx-1`}></i></small>}
						</h4>
					}
					{isLoaded && rows && grid && grid.text2 && uuid == "" && <div className="card-subtitle mb-2 text-muted">{grid.text2}</div>}
					{error && !isLoading && <div className="alert alert-danger" role="alert">{error}</div>}
					{isLoaded && rows && countRows > 0 && uuid == "" &&
						<GridTable rows={rows}
									rowsSelected={rowsSelected}
									rowsEdited={rowsEdited}
									rowsAdded={rowsAdded}
									referencedValuesAdded={referencedValuesAdded}
									referencedValuesRemoved={referencedValuesRemoved}
									columns={grid.columns}
									onSelectRowClick={uuid => this.selectRow(uuid)}
									onEditRowClick={uuid => this.editRow(uuid)}
									onDeleteRowClick={uuid => this.deleteRow(uuid)}
									onAddReferencedValueClick={(fromUuid, col, toGridUuid, uuid, displayString, path) => this.addReferencedValue(fromUuid, col, toGridUuid, uuid, displayString, path)}
									onRemoveReferencedValueClick={(fromUuid, col, toGridUuid, uuid, displayString, path) => this.removeReferencedValue(fromUuid, col, toGridUuid, uuid, displayString, path)}
									inputRef={this.setGridRowRef}
									dbName={dbName}
									token={token} />
					}
					{isLoaded && rows && countRows > 0 && uuid != "" &&
						<GridView row={rows[0]}
									columns={grid.columns}
									referencedValuesAdded={referencedValuesAdded}
									referencedValuesRemoved={referencedValuesRemoved}
									dbName={dbName} />
					}
					<GridFooter isLoading={isLoading}
								grid={grid}
								rows={rows}
								uuid={uuid}
								canAddRows={canAddRows}
								rowsSelected={rowsSelected}
								rowsAdded={rowsAdded}
								rowsEdited={rowsEdited}
								rowsDeleted={rowsDeleted}
								onSelectRowClick={() => this.deselectRows()}
								onAddRowClick={() => this.addRow()}
								onSaveDataClick={() => this.saveData()} />
				</div>
			</div>
		)
	}

	selectRow(selectUuid) { this.setState(state => ({ rowsSelected: [selectUuid] })) }

	deselectRows() { this.setState(state => ({ rowsSelected: [] })) }

	editRow(editUuid) {
		if(!this.state.rowsAdded.includes(editUuid)) {
			this.setState(state => ({
				rowsEdited: state.rowsEdited.filter(uuid => uuid != editUuid).concat(editUuid)
			}))
		}
	}

	addRow() {
		const newRow = { uuid: `${this.props.gridUuid}-${this.state.rows.length+1}` }
		this.setState(state => ({
			rows: state.rows.concat(newRow),
			rowsAdded: state.rowsAdded.concat(newRow.uuid)
		}))
	}

	deleteRow(deleteUuid) {
		if(this.state.rowsAdded.includes(deleteUuid)) {
			this.setState(state => ({
				rowsAdded: state.rowsAdded.filter(uuid => uuid != deleteUuid),
			}))
		} else {
			this.setState(state => ({
				rowsDeleted: state.rowsDeleted.concat(deleteUuid),
			}))
		}
		this.setState(state => ({
			rowsSelected: [],
			rowsEdited: state.rowsEdited.filter(uuid => uuid != deleteUuid),
			rows: state.rows.filter(row => row.uuid != deleteUuid)
		}))
	}

	addReferencedValue(fromUuid, columnName, toGridUuid, uuid, displayString, path) {
		this.setState(state => ({
			referencedValuesAdded: state.referencedValuesAdded.concat({
				fromUuid: fromUuid,
				columnName: columnName,
				toGridUuid: toGridUuid, 
				uuid: uuid,
				displayString: displayString,
				path: path
			}),
			referencedValuesRemoved: state.referencedValuesRemoved.filter(ref => ref.fromUuid != fromUuid || ref.columnName != columnName || ref.uuid != uuid)
		}))
		this.editRow(fromUuid)
	}

	removeReferencedValue(fromUuid, columnName, toGridUuid, uuid, displayString, path) {
		this.setState(state => ({
			referencedValuesRemoved: state.referencedValuesRemoved.concat({
				fromUuid: fromUuid,
				columnName: columnName,
				toGridUuid: toGridUuid, 
				uuid: uuid,
				displayString: displayString,
				path: path
			}),
			referencedValuesAdded: state.referencedValuesAdded.filter(ref => ref.fromUuid != fromUuid || ref.columnName != columnName || ref.uuid != uuid)
		}))
		this.editRow(fromUuid)
	}

	loadData() {
		this.setState({isLoading: true})
		const uri = `/${this.props.dbName}/api/v1/${this.props.gridUuid}${this.props.uuid != "" ? '/' + this.props.uuid : ''}`
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
							grid: result.grid,
							canAddRows: result.canAddRows,
							rows: result.rows,
							error: result.error
						})
					},
					(error) => {
						this.setState({
							isLoading: false,
							isLoaded: false,
							canAddRows: false,							
							rows: [],
							error: error.message
						})
					}
				)
			} else {
				this.setState({
					isLoading: false,
					isLoaded: false,
					canAddRows: false,							
					rows: [],
					error: `[${response.status}] Internal server issue.`
				})
			}
		})
	}

	getInputValues(rows) {
		return rows.map(uuid => {
			const e1 = this.gridInput.get(uuid)
			if(e1) {
				const e2 = Array.from(e1, ([name, value]) => ({ name, value }))
				const e3 = Object.keys(e2).map(key => ({
						col: e2[key].name,
						value: getConvertedValue(e1.get(e2[key].name))
					}))
				const e4 = e3.reduce(
					(hash, {col, value}) => {
						hash[col] = value
						return hash
					},
					{uuid: uuid}
				)
				return e4
			} else {
				return undefined
			}
		})
	}

	saveData() {
		this.setState({isLoading: true})
		const uri = `/${this.props.dbName}/api/v1/${this.props.gridUuid}`
		fetch(uri, {
			method: 'POST',
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
				'Authorization': 'Bearer ' + this.props.token
			},
			body: JSON.stringify({
				rowsAdded: this.getInputValues(this.state.rowsAdded),
				rowsEdited: this.getInputValues(this.state.rowsEdited),
				rowsDeleted: this.getInputValues(this.state.rowsDeleted),
				referencedValuesAdded: this.state.referencedValuesAdded,
				referencedValuesRemoved: this.state.referencedValuesRemoved
			})
		})
		.then(response => {
			const contentType = response.headers.get("content-type")
			if(contentType && contentType.indexOf("application/json") !== -1) {
				return response.json().then(	
					(result) => {
						this.setState({
							isLoading: false,
							isLoaded: true,
							grid: result.grid,
							rows: result.rows,
							error: result.error,
							rowsEdited: [],
							rowsSelected: [],
							rowsAdded: [],
							rowsDeleted: [],
							referencedValuesAdded: [],
							referencedValuesRemoved: []
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
	}
}

Grid.defaultProps = {
	gridUuid: '',
	uuid: ''
}

class GridFooter extends React.Component {
	render() {
		const { isLoading, grid, rows, uuid, canAddRows, rowsEdited, rowsAdded, rowsDeleted } = this.props
		const countRows = rows ? rows.length : 0
		const countRowsAdded = rowsAdded ? rowsAdded.length : 0
		const countRowsEdited = rowsEdited ? rowsEdited.length : 0
		const countRowsDeleted = rowsDeleted ? rowsDeleted.length : 0
		return (
			<div onClick={() => this.props.onSelectRowClick()}>
				{countRows == 0 && <small className="text-muted px-1">No data</small>}
				{countRows == 1 && <small className="text-muted px-1">{countRows} row</small>}
				{countRows > 1 && <small className="text-muted px-1">{countRows} rows</small>}
				{countRowsAdded > 0 && <small className="text-muted px-1">({countRowsAdded} added)</small>}
				{countRowsEdited > 0 && <small className="text-muted px-1">({countRowsEdited} edited)</small>}
				{countRowsDeleted > 0 && <small className="text-muted px-1">({countRowsDeleted} deleted)</small>}
				{!isLoading && grid && <a href={grid.path}><i className="bi bi-box-arrow-up-right mx-1"></i></a>}
				{!isLoading && grid && uuid == "" && canAddRows &&
					<button type="button" className="btn btn-outline-success btn-sm mx-1"
							onClick={this.props.onAddRowClick}>
						Add <i className="bi bi-plus-circle"></i>
					</button>
				}
				{!isLoading && countRowsAdded + countRowsEdited + countRowsDeleted > 0 &&
					<button type="button" className="btn btn-outline-primary btn-sm mx-1"
							onClick={this.props.onSaveDataClick}>
						Save <i className="bi bi-save"></i>
					</button>
				}
				{isLoading && <Spinner />}
			</div>
		)
	}
}

function getColumnType(type) {
	switch(type) {
		case UuidIntColumnType:
			return "number"
		case UuidPasswordColumnType:
			return "password"
		case UuidReferenceColumnType:
			return "reference"
		default:
			return "text"
	}
}

function getColumnValuesForRow(columns, row, withTimeStamps) {
	const cols = []
	{columns && columns.map(
		column => {
			let type = getColumnType(column.typeUuid)
			if(type == "reference") {
				cols.push({
					name: column.name,
					label: column.label,
					values: getColumnValueForReferencedRow(column, row),
					typeUuid: column.typeUuid,
					gridPromptUuid: column.gridPromptUuid,
					type: type,
					readonly: false
				})
			} else {
				cols.push({
					name: column.name,
					label: column.label,
					value: row[column.name],
					typeUuid: column.typeUuid,
					type: type,
					readonly: false
				})	
			}
		}
	)}
	if(withTimeStamps) {
		cols.push({name: "uuid", label: "Identifier", value: row.uuid, typeUuid: UuidUuidColumnType, type: "text", readonly: true})
		cols.push({name: "revision", label: "Revision", value: row.revision, type: "number", readonly: true})
	}
	return cols
}

function getColumnValueForReferencedRow(column, row) {
	let output = []
	if(row.references) {
		row.references.map(
			ref => {
				if(ref.name == column.name && ref.rows) {
					ref.rows.map(
						refRow => output.push({
							uuid: refRow.uuid,
							displayString: refRow.displayString,
							path: refRow.path
						})
					)
				}
			}
		)
	}
	return output
}

function getCellValue(type, value) {
	if(type == 'password') return '*****'
	return value
}

function getConvertedValue(cell) {
	switch(cell.type) {
		case "number":
			return Number(cell.value)
		case "text":
			return String(cell.value)
		default:
			return cell.value
	}
}

class Spinner extends React.Component {
	render() {
		return (
			<span className="spinner-grow spinner-grow-sm ms-1" role="status">
				<span className="visually-hidden">Loading...</span>
			</span>
		)
	}
}
