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
		const countRowsAdded = rowsAdded.length
		const countRowsEdited = rowsEdited.length
		const rowsSelectedAndEdited = rowsSelected.filter(uuid => !rowsEdited.includes(uuid))
		const countRowsSelectedAndEdited = rowsSelectedAndEdited.length
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
					</h6>
					{error && !isLoading && !isLoaded && <div className="alert alert-danger" role="alert">{error}</div>}
					{error && !isLoading && isLoaded && <div className="alert alert-primary" role="alert">{error}</div>}

					{isLoaded && rows && countRows > 0 && this.props.uuid == "" &&
						 <GridTable rows={rows}
						 			rowsSelected={rowsSelected}
									rowsEdited={rowsEdited}
									rowsAdded={rowsAdded}
						 			onRowClick={row => this.toggleSelection(row)}/>
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
					{isLoaded && rows && this.props.uuid == "" && countRowsSelected > 0 &&
						<button
							type="button"
							className="btn btn-outline-danger btn-sm mx-1"
							onClick={() => this.deleteRows()}>
							Delete selected <i className="bi bi-dash-circle"></i>
						</button>
					}
					{isLoaded && rows && this.props.uuid == "" && countRowsSelectedAndEdited > 0 &&
						<button
							type="button"
							className="btn btn-outline-secondary btn-sm mx-1"
							onClick={() => this.setRowsEdited()}>
							Edit selected <i className="bi bi-pencil"></i>
						</button>
					}
					{isLoaded && rows && this.props.uuid == "" && (countRowsAdded > 0 || countRowsEdited > 0) &&
						<button
							type="button"
							className="btn btn-outline-primary btn-sm mx-1"
							onClick={() => this.deleteRows()}>
							Save changes <i className="bi bi-save"></i>
						</button>
					}
				</div>
			</div>
		)
	}

	toggleSelection(row) {
		if(!this.isRowSelected(row)) {
			this.setState(state => ({
				rowsSelected: state.rowsSelected.concat(row.uuid),
			}))
		}
		else {
			this.setState(state => ({
				rowsSelected: state.rowsSelected.filter(uuid => uuid != row.uuid),
			}))
		}
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

	addRow() {
		const newRow = {
			uri: "new row",
			uuid: `${this.state.rows.length+1}`
		}
		this.setState(state => ({
			rows: state.rows.concat(newRow),
			rowsEdited: state.rowsEdited.concat(newRow.uuid),
			rowsSelected: state.rowsSelected.concat(newRow.uuid),
			rowsAdded: state.rowsAdded.concat(newRow.uuid),
		}))
	}

	deleteRows() {
		this.setState(state => ({
			rows: state.rows.filter(row => !this.isRowSelected(row)),
			rowsSelected: []
		}))
	}

	setRowsEdited() {
		this.setState(state => ({
			rowsEdited: state.rowsEdited.concat(state.rowsSelected),
		}))
	}
}

Grid.defaultProps = {
	gridUri: '',
	uuid: ''
}
