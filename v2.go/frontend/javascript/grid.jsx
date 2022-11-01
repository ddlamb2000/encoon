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
			items: [],
			count: 0,
			itemsSelected: [],
			countSelected: 0,
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
							items: result.items,
							count: result.count,
							error: result.error
						})
					},
					(error) => {
						this.setState({
							isLoading: false,
							isLoaded: false,
							items: [],
							error: error.message
						})
					}
				)
			} else {
				this.setState({
					isLoading: false,
					isLoaded: false,
					items: [],
					error: `[${response.status}] Internal server issue.`
				})
			}
		})
	}

	render() {
		const { isLoading, isLoaded, error, grid, items, count, itemsSelected, countSelected } = this.state
		return (
			<div className="card mt-2 mb-2">
				<div className="card-body">
					{isLoading && <Spinner />}
					{grid && <h4 className="card-title">{grid.text01}</h4>}
					{error && !isLoading && !isLoaded && <div className="alert alert-danger" role="alert">{error}</div>}
					{error && !isLoading && isLoaded && <div className="alert alert-primary" role="alert">{error}</div>}
					{isLoaded && items && count == 0 && <div className="alert alert-secondary" role="alert">No data</div>}

					{isLoaded && items && count > 0 && this.props.uuid == "" &&
						 <TableRows items={items}
						 			itemsSelected={itemsSelected}
						 			onRowItemClick={item => this.toggleSelection(item)}/>
					}

					{isLoaded && items && count > 0 && this.props.uuid != "" && <TableSingleRow item={items[0]} />}

					{isLoaded && items && this.props.uuid == "" && count == 1 && <small className="text-muted px-2">{count} row</small>}
					{isLoaded && items && this.props.uuid == "" && count > 1 && <small className="text-muted px-2">{count} rows</small>}

					{isLoaded && items &&
						<button
							type="button"
							className="btn btn-light btn-sm px-2"
							onClick={() => this.addItem()}>
							Add <img src="/icons/plus-circle.svg" role="img" alt="Add row"></img>
						</button>
					}

					{isLoaded && items && this.props.uuid == "" && countSelected > 0 &&
						<small className="text-muted px-2">{countSelected} selected</small>
					}
					{isLoaded && items && this.props.uuid == "" && countSelected > 0 &&
						<button
							type="button"
							className="btn btn-light btn-sm px-2"
							onClick={() => this.deleteItems()}>
							Delete <img src="/icons/dash-circle.svg" role="img" alt="Delete row"></img>
						</button>
					}

				</div>
			</div>
		)
	}

	toggleSelection(item) {
		if(!this.isItemSelected(item)) {
			this.setState(state => ({
				itemsSelected: state.itemsSelected.concat(item.uuid),
				countSelected: state.countSelected + 1
			}))
		}
		else {
			this.setState(state => ({
				itemsSelected: state.itemsSelected.filter(uuid => uuid != item.uuid),
				countSelected: state.countSelected - 1
			}))
		}
	}

	isItemSelected(item) {
		return this.state.itemsSelected.some(uuid => uuid == item.uuid)
	}

	addItem() {
		this.setState(state => ({
			items: state.items.concat({
				uri: "????",
				uuid: `NEW-${this.state.count+1}`,
				editable: true,
				added: true
			}),
			count: state.count + 1
		}))
	}

	deleteItems() {
		this.setState(state => ({
			items: state.items.filter(item => !this.isItemSelected(item)),
			itemsSelected: [],
		}))
		this.setState(state => ({
			count: state.items.length,
			countSelected: 0
		}))
	}
}

Grid.defaultProps = {
	gridUri: '',
	uuid: ''
}

class TableRows extends React.Component {
	render() {
		return (
			<table className="table table-hover table-sm">
				<thead className="table-light">
					<tr>
						<th scope="col"><span>&nbsp;</span></th>
						<th scope="col">Uri</th>
						<th scope="col">Text01</th>
						<th scope="col">Text02</th>
						<th scope="col">Text03</th>
						<th scope="col">Text04</th>
						<th scope="col">Uuid</th>
					</tr>
				</thead>
				<tbody>
					{this.props.items.map(
						item => <TableRowItemSingleLine key={item.uuid} 
														item={item}
														itemSelected={this.isItemSelected(item)}
														onRowItemClick={item => this.props.onRowItemClick(item)} />
					)}
				</tbody>
			</table>
		)
	}

	isItemSelected(item) {
		return this.props.itemsSelected.some(uuid => uuid == item.uuid)
	}
}

class TableRowItemSingleLine extends React.Component {
	render() {
		const variant = (this.props.itemSelected ? "table-warning" : "")
		return (
			<tr className={variant} onClick={() => this.props.onRowItemClick(this.props.item)}>
				<td>
					{this.props.item.added && <span>*</span>}
					{!this.props.item.added && <span>&nbsp;</span>}
				</td>
				{this.props.item.editable && <td><input></input></td>}
				{!this.props.item.editable && <td>{this.props.item.uri}</td>}
				<td>{this.props.item.text01}</td>
				<td>{this.props.item.text02}</td>
				<td>{this.props.item.text03}</td>
				<td>{this.props.item.text04}</td>
				<td scope="row"><a href={this.props.item.path}>{this.props.item.uuid}</a></td>
			</tr>
		)
	}
}

class TableSingleRow extends React.Component {
	render() {
		return (
			<div>
				<h6 className="card-subtitle mb-2 text-muted">{this.props.item.uuid}</h6>
				<table className="table table-hover table-sm">
					<thead className="table-light"></thead>
					<tbody>
						<tr><td>Uri</td><td>{this.props.item.uri}</td></tr>
						<tr><td>Text01</td><td>{this.props.item.text01}</td></tr>
						<tr><td>Text02</td><td>{this.props.item.text02}</td></tr>
						<tr><td>Text03</td><td>{this.props.item.text03}</td></tr>
						<tr><td>Text04</td><td>{this.props.item.text04}</td></tr>
					</tbody>
				</table>
			</div>
		)
	}
}
