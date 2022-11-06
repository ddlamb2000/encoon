// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Grid extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: false,
			isLoaded: false,
			isLoading: false,
			grid: [],
			rows: [],
			rowsSelected: [],
			rowsEdited: [],
			rowsAdded: []
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
		this.setState({isLoading: true})
		const uri = `/${this.props.dbName}/api/v1/${this.props.gridUri}${this.props.uuid != "" ? '/' + this.props.uuid : ''}`
		fetch(uri, {
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
				'Authorization': 'Bearer ' + this.props.token
			}
		})
		.then(response => {
			const contentType = response.headers.get("content-type");
			if(contentType && contentType.indexOf("application/json") !== -1) {
				return response.json().then(	
					(result) => {
						this.setState({
							isLoading: false,
							isLoaded: true,
							grid: result.grid,
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
	}

	render() {
		const { isLoading, isLoaded, error, grid, rows, rowsSelected, rowsEdited, rowsAdded } = this.state
		const countRows = rows ? rows.length : 0
		const countRowsSelected = rowsSelected.length
		const countRowsEdited = rowsEdited.length
		return (
			<div className="card mt-2 mb-2">
				<div className="card-body">
					{isLoading && <Spinner />}
					{grid && <h4 className="card-title">{grid.text01}</h4>}
					<h6 className="card-subtitle mb-2 text-muted">
						{isLoaded && rows && this.props.uuid == "" && countRows == 1 && <small className="text-muted px-2">{countRows} row</small>}
						{isLoaded && rows && this.props.uuid == "" && countRows > 1 && <small className="text-muted px-2">{countRows} rows</small>}
						{isLoaded && rows && this.props.uuid == "" && countRowsSelected > 0 &&
							<small className="text-muted px-2">({countRowsSelected} selected)</small>
						}
						{isLoaded && rows && countRows == 0 && <small className="text-muted px-2">No data</small>}
						{grid &&<a href={grid.path}><i className="bi bi-box-arrow-up-right ms-2"></i></a>}
					</h6>
					{error && !isLoading && !isLoaded && <div className="alert alert-danger" role="alert">{error}</div>}
					{error && !isLoading && isLoaded && <div className="alert alert-primary" role="alert">{error}</div>}

					{isLoaded && rows && countRows > 0 && this.props.uuid == "" &&
						 <GridTable rows={rows}
						 			rowsSelected={rowsSelected}
									rowsEdited={rowsEdited}
									rowsAdded={rowsAdded}
						 			onRowClick={uuid => this.selectRow(uuid)}
						 			onEditRowClick={uuid => this.editRow(uuid)}
						 			onDeleteRowClick={uuid => this.deleteRow(uuid)}
									inputRef={this.setGridRowRef} />
					}

					{isLoaded && rows && countRows > 0 && this.props.uuid != "" && <GridView row={rows[0]} />}

					{isLoaded && rows &&
						<button
							type="button"
							className="btn btn-outline-success btn-sm mx-1"
							onClick={() => this.addRow()}>
							Add row <i className="bi bi-plus-circle"></i>
						</button>
					}
					{isLoaded && rows && this.props.uuid == "" && countRowsEdited > 0 &&
						<button
							type="button"
							className="btn btn-outline-primary btn-sm mx-1"
							onClick={() => this.saveData()}>
							Save changes <i className="bi bi-save"></i>
						</button>
					}
				</div>
			</div>
		)
	}

	selectRow(selectUuid) {
		this.setState(state => ({
			rowsSelected: state.rowsSelected.filter(uuid => uuid == selectUuid).concat(selectUuid)
		}))
	}

	editRow(editUuid) {
		this.setState(state => ({
			rowsEdited: state.rowsEdited.filter(uuid => uuid != editUuid).concat(editUuid)
		}))
	}

	addRow() {
		const newRow = { uuid: `${this.props.gridUri}-${this.state.rows.length+1}` }
		this.setState(state => ({
			rows: state.rows.concat(newRow),
			rowsAdded: state.rowsAdded.concat(newRow.uuid),
			rowsEdited: state.rowsEdited.concat(newRow.uuid)
		}))
	}

	deleteRow(deleteUuid) {
		this.setState(state => ({
			rows: state.rows.filter(row => row.uuid != deleteUuid)
		}))
	}

	setRowsEdited() {
		this.setState(state => ({
			rowsEdited: state.rowsEdited.concat(state.rowsSelected),
		}))
	}

	isRowSelected(row) {
		return this.state.rowsSelected.includes(row.uuid)
	}

	isRowEdited(row) {
		return this.state.rowsEdited.includes(row.uuid)
	}

	isRowAdded(row) {
		return this.state.rowsAdded.includes(row.uuid)
	}

	getInputValues(rows) {
		return rows.map(uuid => {
			const inputElement = this.gridInput.get(uuid)
			const e0 = Array.from(inputElement, ([name, value]) => ({ name, value }))
			const e1 = Object.keys(e0).map(key => 
				({
					col: e0[key].name,
					value: inputElement.get(e0[key].name).value
				}))
			const e2 = e1.reduce(
				(hash, {col, value}) => {
					hash[col] = value
					return hash
				},
				{uuid: uuid}
			)
			return e2
		})
	}

	saveData() {
		const rowsEditedAndNotAdded = this.state.rowsEdited.filter(uuid => !this.state.rowsAdded.includes(uuid))
		const rowsAdded = this.getInputValues(this.state.rowsAdded)
		const rowsEdited = this.getInputValues(rowsEditedAndNotAdded)
		const uri = `/${this.props.dbName}/api/v1/${this.props.gridUri}`
		fetch(uri, {
			method: 'POST',
			headers: {
				'Accept': 'application/json',
				'Content-Type': 'application/json',
				'Authorization': 'Bearer ' + this.props.token
			},
			body: JSON.stringify({ rowsAdded: rowsAdded, rowsEdited: rowsEdited })
		})
		.then(response => {
			const contentType = response.headers.get("content-type");
			if(contentType && contentType.indexOf("application/json") !== -1) {
				return response.json().then(	
					(result) => {
						if(response.status == 200) {
							this.setState({
								rowsEdited: [],
								rowsSelected: [],
								rowsAdded: [],
							})
						}
						else {
							alert(result.error)
							this.setState({
								error: result.error
							})
						}
					},
					(error) => {
						this.setState({
							error: error.message
						})
					}
				)
			} else {
				this.setState({
					error: `[${response.status}] Internal server issue.`
				})
			}
		})
	}
}

Grid.defaultProps = {
	gridUri: '',
	uuid: ''
}
