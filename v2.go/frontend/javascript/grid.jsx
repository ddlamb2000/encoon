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
			rowsAdded: [],
			rowsDeleted: []
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
		const { isLoading, isLoaded, error, grid, rows, rowsSelected, rowsEdited, rowsAdded, rowsDeleted } = this.state
		const countRows = rows ? rows.length : 0
		const countRowsAdded = rowsAdded.length
		const countRowsEdited = rowsEdited.length
		const countRowsDeleted = rowsDeleted.length
		return (
			<div className="card mt-2 mb-2">
				<div className="card-body">
					{grid && <h4 className="card-title">{grid.text02}</h4>}
					{isLoaded && rows && grid && grid.text03 && <div className="card-subtitle mb-2 text-muted">{grid.text03}</div>}
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
					{error && !isLoading && !isLoaded && <div className="alert alert-danger" role="alert">{error}</div>}
					{error && !isLoading && isLoaded && <div className="alert alert-primary" role="alert">{error}</div>}
					{isLoaded && rows && countRows == 0 && <small className="text-muted px-1">No data</small>}
					{isLoaded && rows && grid && <a href={grid.path}><i className="bi bi-box-arrow-up-right mx-1"></i></a>}
					{isLoaded && rows &&
						<button
							type="button"
							className="btn btn-outline-success btn-sm mx-1"
							onClick={() => this.addRow()}>
							Add row <i className="bi bi-plus-circle"></i>
						</button>
					}
					{isLoaded && rows && this.props.uuid == "" && countRowsAdded + countRowsEdited + countRowsDeleted > 0 &&
						<button
							type="button"
							className="btn btn-outline-primary btn-sm mx-1"
							onClick={() => this.saveData()}>
							Save changes <i className="bi bi-save"></i>
						</button>
					}
					{isLoaded && rows && countRows == 1 && <small className="text-muted px-1">{countRows} row</small>}
					{isLoaded && rows && countRows > 1 && <small className="text-muted px-1">{countRows} rows</small>}
					{isLoaded && rows && countRowsAdded > 0 && <small className="text-muted px-1">({countRowsAdded} added)</small>}
					{isLoaded && rows && countRowsEdited > 0 && <small className="text-muted px-1">({countRowsEdited} edited)</small>}
					{isLoaded && rows && countRowsDeleted > 0 && <small className="text-muted px-1">({countRowsDeleted} deleted)</small>}
					{isLoading && <Spinner />}
				</div>
			</div>
		)
	}

	selectRow(selectUuid) {
		this.setState(state => ({
			rowsSelected: [selectUuid]
		}))
	}

	editRow(editUuid) {
		if(!this.state.rowsAdded.includes(editUuid)) {
			this.setState(state => ({
				rowsEdited: state.rowsEdited.filter(uuid => uuid != editUuid).concat(editUuid)
			}))
		}
	}

	addRow() {
		const newRow = { uuid: `${this.props.gridUri}-${this.state.rows.length+1}` }
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

	getInputValues(rows) {
		return rows.map(uuid => {
			const e1 = this.gridInput.get(uuid)
			const e2 = Array.from(e1, ([name, value]) => ({ name, value }))
			const e3 = Object.keys(e2).map(key => ({
					col: e2[key].name,
					value: e1.get(e2[key].name).value
				}))
			const e4 = e3.reduce(
				(hash, {col, value}) => {
					hash[col] = value
					return hash
				},
				{uuid: uuid}
			)
			return e4
		})
	}

	loadData() {
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

	saveData() {
		this.setState({isLoading: true})
		const uri = `/${this.props.dbName}/api/v1/${this.props.gridUri}`
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
				rowsDeleted: this.getInputValues(this.state.rowsDeleted)
			})
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
							error: result.error,
							rowsEdited: [],
							rowsSelected: [],
							rowsAdded: [],
							rowsDeleted: []
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
	gridUri: '',
	uuid: ''
}
