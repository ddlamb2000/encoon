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
		this.inputMap = new Map()
		this.richTextMap = new Map()
		this.setGridRowRef = element => {
			if(element != undefined) {
				const uuid = element.getAttribute("uuid")
				const columnName = element.getAttribute("column")
				if(uuid != "" && columnName != "") {
					const inputRowMap = this.inputMap.get(uuid)
					if(inputRowMap) {
						inputRowMap.set(columnName, element)
					}
					else {
						const inputRowMap = new Map()
						inputRowMap.set(columnName, element)
						this.inputMap.set(uuid, inputRowMap)
					}
				}
			}
		}
	}

	componentDidMount() {
		this.loadData()
	}

	componentDidUpdate(prevProps) {
		if(trace) console.log("[Grid.componentDidUpdate()] **** this.props.gridUuid=", this.props.gridUuid, ", this.props.uuid=", this.props.uuid)
		if(this.props.gridUuid != prevProps.gridUuid || this.props.uuid != prevProps.uuid) {
			this.loadData()
		}
	}

	render() {
		const { isLoading, 
				isLoaded, 
				error, 
				grid, 
				canAddRows, 
				rows, 
				rowsSelected, 
				rowsEdited, 
				rowsAdded, 
				rowsDeleted, 
				referencedValuesAdded, 
				referencedValuesRemoved } = this.state
		const { dbName, 
				token, 
				gridUuid, 
				uuid,
				filterColumnOwned,
				filterColumnName,
				filterColumnLabel,
				filterColumnGridUuid,
				filterColumnDisplayString,
				innerGrid,
				miniGrid,
				gridTitle,
				gridSubTitle,
			 	noEdit } = this.props
		const countRows = rows ? rows.length : 0
		const columns = miniGrid ? (grid && grid.columns != undefined ? grid.columns.slice(0,1) : []) : grid.columns
		if(trace) console.log("[Grid.render()] gridUuid=", gridUuid, ", uuid=", uuid)
		return (
			<div className={!innerGrid ? "card my-4" : ""}>
				<div className={!innerGrid ? "card-body" : ""}>
					{isLoaded && rows && grid && uuid == "" && !innerGrid &&
						<h5 className="card-title">
							{grid.text1} {grid.text3 && <small><i className={`bi bi-${grid.text3} mx-1`}></i></small>}
						</h5>
					}
					{isLoaded && filterColumnLabel && filterColumnName && !innerGrid &&
						<div className="card-subtitle mb-2"><mark>{filterColumnLabel} = {filterColumnDisplayString}</mark></div>
					}
					{isLoaded && rows && grid && grid.text2 && uuid == "" && !innerGrid && 
						<div className="card-subtitle mb-2 text-muted">{grid.text2}</div>
					}
					{isLoaded && rows && countRows > 0 && grid && uuid == "" && innerGrid && gridTitle &&
						<h5 className="card-title">{gridTitle}</h5>
					}
					{isLoaded && rows && countRows > 0 && grid && uuid == "" && innerGrid && gridSubTitle &&
						<div className="card-subtitle mb-2 text-muted">{gridSubTitle}  <small><i className={`bi bi-grid-3x3 mx-1`}></i></small></div>
					}
					{error && !isLoading && <div className="alert alert-danger" role="alert">{error}</div>}
					{isLoaded && rows && countRows > 0 && uuid == "" &&
						<GridTable rows={rows}
									columns={columns}
									grid={grid}
									rowsSelected={rowsSelected}
									rowsEdited={rowsEdited}
									rowsAdded={rowsAdded}
									referencedValuesAdded={referencedValuesAdded}
									referencedValuesRemoved={referencedValuesRemoved}
									onSelectRowClick={uuid => this.selectRow(uuid)}
									onEditRowClick={uuid => this.editRow(uuid)}
									onDeleteRowClick={uuid => this.deleteRow(uuid)}
									onAddReferencedValueClick={reference => this.addReferencedValue(reference)}
									onRemoveReferencedValueClick={reference => this.removeReferencedValue(reference)}
									inputRef={this.setGridRowRef}
									dbName={dbName}
									token={token}
									navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
									filterColumnOwned={filterColumnOwned}
									filterColumnName={filterColumnName}
									filterColumnGridUuid={filterColumnGridUuid}
									createRichTextField={(id, value) => this.createRichTextField(id, value)}
									deleteRichTextField={id => this.deleteRichTextField(id)}
									miniGrid={miniGrid} />
					}
					{isLoaded && rows && countRows > 0 && uuid != "" &&
						<GridView row={rows[0]}
									columns={grid.columns}
									columnsUsage={grid.columnsUsage}
									grid={grid}
									rowsSelected={rowsSelected}
									rowsEdited={rowsEdited}
									rowsAdded={rowsAdded}
									referencedValuesAdded={referencedValuesAdded}
									referencedValuesRemoved={referencedValuesRemoved}
									onSelectRowClick={uuid => this.selectRow(uuid)}
									onEditRowClick={uuid => this.editRow(uuid)}
									onDeleteRowClick={uuid => this.deleteRow(uuid)}
									onAddReferencedValueClick={reference => this.addReferencedValue(reference)}
									onRemoveReferencedValueClick={reference => this.removeReferencedValue(reference)}
									inputRef={this.setGridRowRef}
									dbName={dbName}
									navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
									token={this.props.token}
									loadParentData={() => this.loadData()}
									createRichTextField={(id, value) => this.createRichTextField(id, value)}
									deleteRichTextField={id => this.deleteRichTextField(id)} />
					}
					{!noEdit &&
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
									onSaveDataClick={() => this.saveData()}
									navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
									miniGrid={this.props.miniGrid} />
					}
				</div>
			</div>
		)
	}

	selectRow(selectUuid) {
		if(trace) console.log("[Grid.selectRow()] ", selectUuid)
		this.setState(state => ({ rowsSelected: [selectUuid] })) 
	}

	deselectRows() {
		if(trace) console.log("[Grid.deselectRows()] ")
		this.setState(state => ({ rowsSelected: [] })) 
	}

	editRow(editUuid) {
		if(trace) console.log("[Grid.editRow()] ", editUuid)
		if(!this.state.rowsAdded.includes(editUuid)) {
			this.setState(state => ({
				rowsEdited: state.rowsEdited.filter(uuid => uuid != editUuid).concat(editUuid)
			}))
		}
	}

	addRow() {
		if(trace) console.log("[Grid.addRow()] ")
		const newRow = { uuid: `${this.props.gridUuid}-${this.state.rows.length+1}` }
		this.setState(state => ({
			rows: state.rows.concat(newRow),
			rowsAdded: state.rowsAdded.concat(newRow.uuid)
		}))
	}

	deleteRow(deleteUuid) {
		if(trace) console.log("[Grid.deleteRow()] ", deleteUuid)
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

	addReferencedValue(reference) {
		if(trace) console.log("[Grid.addReferencedValue()] ", reference)
		this.setState(state => ({
			referencedValuesAdded: state.referencedValuesAdded.concat({
				fromUuid: reference.fromUuid,
				columnUuid: reference.columnUuid,
				owned: reference.owned,
				columnName: reference.columnName,
				toGridUuid: reference.toGridUuid, 
				uuid: reference.uuid,
				displayString: reference.displayString
			}),
			referencedValuesRemoved: state.referencedValuesRemoved.
				filter(ref => 
						ref.fromUuid != reference.fromUuid || 
						ref.columnUuid != reference.columnUuid || 
						ref.uuid != reference.uuid)
		}))
		this.editRow(reference.fromUuid)
	}

	removeReferencedValue(reference) {
		if(trace) console.log("[Grid.removeReferencedValue()] ", reference)
		this.setState(state => ({
			referencedValuesRemoved: state.referencedValuesRemoved.concat({
				fromUuid: reference.fromUuid,
				columnUuid: reference.columnUuid,
				owned: reference.owned,
				columnName: reference.columnName,
				toGridUuid: reference.toGridUuid, 
				uuid: reference.uuid,
				displayString: reference.displayString
			}),
			referencedValuesAdded: state.referencedValuesAdded.
				filter(ref => 
						ref.fromUuid != reference.fromUuid || 
						ref.columnUuid != reference.columnUuid || 
						ref.uuid != reference.uuid)
		}))
		this.editRow(reference.fromUuid)
	}

	createRichTextField(id, value) {
		const richText = new Quill('#' + id, {
			modules: {
				toolbar: [
					[{ header: [1, 2, 3, false] }],
					['bold', 'italic', 'underline', 'strike','script', 'code'],
					['blockquote', 'list', 'align', 'code-block']
				]
			},
			theme: 'snow'
		})
		this.setRichTextValue(richText, value)
		this.richTextMap.set(id, richText)
	}

	setRichTextValue(richText, value) {
		if(value) {
			try {
				console.log(value)
				richText.setContents(JSON.parse(value))
			} catch (error) {
				console.error("Invalid value", value)
			}
		}
	}

	deleteRichTextField(id) {
		const richText = this.richTextMap.get(id)
		if(richText) {
			this.richTextMap.delete(id)
		}
	}

	loadData() {
		this.setState({
			error: "",
			isLoaded: false,
			isLoading: true,
			grid: [],
			canAddRows: false,
			rows: [],
			rowsSelected: [],
			rowsEdited: [],
			rowsAdded: [],
			rowsDeleted: [],
			referencedValuesAdded: [],
			referencedValuesRemoved: []
		})
		const { dbName, token, gridUuid, uuid, filterColumnOwned, filterColumnName, filterColumnGridUuid, filterColumnValue } = this.props
		const uuidFilter = uuid != "" ? '/' + uuid : ''
		const columnFilter = filterColumnName && filterColumnGridUuid && filterColumnValue ?
								'?filterColumnOwned=' + filterColumnOwned +
								'&filterColumnName=' + filterColumnName + 
								'&filterColumnGridUuid=' + filterColumnGridUuid + 
								'&filterColumnValue=' + filterColumnValue : ''
		const uri = `/${dbName}/api/v1/${gridUuid}${uuidFilter}${columnFilter}`
		fetch(uri, {
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
				'Authorization': 'Bearer ' + token
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
								grid: result.response.grid,
								canAddRows: result.response.canAddRows,
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
			const e1 = this.inputMap.get(uuid)
			if(e1) {
				const e2 = Array.from(e1, ([name, value]) => ({ name, value }))
				const e3 = Object.keys(e2).map(key => ({
						col: e2[key].name,
						value: this.getConvertedValue(e1.get(e2[key].name))
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

	getConvertedValue(cell) {
		switch(cell.getAttribute('typeuuid')) {
			case UuidTextColumnType: return String(cell.value)
			case UuidIntColumnType: return Number(cell.value)
			case UuidUuidColumnType: return String(cell.value)
			case UuidBooleanColumnType: return String(cell.checked)
			case UuidRichTextColumnType:
				const richText = this.richTextMap.get(cell.id)
				if(richText != undefined) {
					const content = richText.getContents()
					const contentString = JSON.stringify(content)
					return contentString
				}				
		}
		return cell.value
	}
	
	saveData() {
		const { dbName, token, gridUuid, uuid, filterColumnOwned, filterColumnName, filterColumnGridUuid, filterColumnValue } = this.props
		this.setState({isLoading: true})
		const uuidFilter = uuid != "" ? '/' + uuid : ''
		const columnFilter = filterColumnName && filterColumnGridUuid && filterColumnValue ? 
								'?filterColumnOwned=' + filterColumnOwned +
								'&filterColumnName=' + filterColumnName + 
								'&filterColumnGridUuid=' + filterColumnGridUuid +
								'&filterColumnValue=' + filterColumnValue : ''
		const uri = `/${dbName}/api/v1/${gridUuid}${uuidFilter}${columnFilter}`
		fetch(uri, {
			method: 'POST',
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
				'Authorization': 'Bearer ' + token
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
			if(contentType && contentType.indexOf("application/json") != -1) {
				return response.json().then(	
					(result) => {
						if(result.response != undefined) {
							this.setState({
								isLoading: false,
								isLoaded: true,
								grid: result.response.grid,
								rows: result.response.rows,
								error: result.response.error,
								rowsEdited: [],
								rowsSelected: [],
								rowsAdded: [],
								rowsDeleted: [],
								referencedValuesAdded: [],
								referencedValuesRemoved: []
							})
							if(this.props.loadParentData != undefined) this.props.loadParentData()
						} else {
							this.setState({
								isLoading: false,
								isLoaded: true,
								error: result.error,
								rowsEdited: [],
								rowsSelected: [],
								rowsAdded: [],
								rowsDeleted: [],
								referencedValuesAdded: [],
								referencedValuesRemoved: []
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
	}
}

Grid.defaultProps = {
	gridUuid: '',
	uuid: ''
}

class GridFooter extends React.Component {
	render() {
		const { grid, rows, uuid, canAddRows, rowsEdited, rowsAdded, rowsDeleted, isLoading, miniGrid } = this.props
		const countRows = rows ? rows.length : 0
		const countRowsAdded = rowsAdded ? rowsAdded.length : 0
		const countRowsEdited = rowsEdited ? rowsEdited.length : 0
		const countRowsDeleted = rowsDeleted ? rowsDeleted.length : 0
		return (
			<nav className='mb-3' onClick={() => this.props.onSelectRowClick()}>
				{isLoading && <Spinner />}
				{!isLoading && !miniGrid && countRows == 0 && <small className="text-muted px-1">No data</small>}
				{!isLoading && !miniGrid && countRows == 1 && uuid == '' && <small className="text-muted px-1">{countRows} row</small>}
				{!isLoading && !miniGrid && countRows > 1 && <small className="text-muted px-1">{countRows} rows</small>}
				{!isLoading && !miniGrid && countRowsAdded > 0 && <small className="text-muted px-1">({countRowsAdded} added)</small>}
				{!isLoading && !miniGrid && countRowsEdited > 0 && <small className="text-muted px-1">({countRowsEdited} edited)</small>}
				{!isLoading && !miniGrid && countRowsDeleted > 0 && <small className="text-muted px-1">({countRowsDeleted} deleted)</small>}
				{!isLoading && grid && uuid == "" && canAddRows &&
					<button type="button"
							className="btn btn-outline-success btn-sm mx-1"
							onClick={this.props.onAddRowClick}>
						Add <i className="bi bi-plus-circle"></i>
					</button>
				}
				{!isLoading && countRowsAdded + countRowsEdited + countRowsDeleted > 0 &&
					<button type="button"
							className="btn btn-outline-primary btn-sm mx-1"
							onClick={this.props.onSaveDataClick}>
						Save <i className="bi bi-save"></i>
					</button>
				}
			</nav>
		)
	}
}

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
