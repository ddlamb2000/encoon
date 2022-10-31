// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Grid extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: false,
			isLoaded: false,
			isLoading: false,
			items: [],
			count: 0,
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
		const { isLoading, isLoaded, error, items, count } = this.state
		return (
			<div className="card mt-2 mb-2">
				<div className="card-body">
					<h4 className="card-title">{this.props.gridUri}{isLoading && <Spinner />}</h4>
					{error && !isLoading && !isLoaded && <div className="alert alert-danger" role="alert">{error}</div>}
					{error && !isLoading && isLoaded && <div className="alert alert-primary" role="alert">{error}</div>}
					{isLoaded && items && count == 0 && <div className="alert alert-secondary" role="alert">No data</div>}
					{isLoaded && items && count > 0 && this.props.uuid == "" && <TableRows items={items} />}
					{isLoaded && items && count > 0 && this.props.uuid != "" && <TableSingleRow item={items[0]} />}
					{isLoaded && items && this.props.uuid == "" && count == 1 && <p><small className="text-muted">{count} row</small></p>}
					{isLoaded && items && this.props.uuid == "" && count > 1 && <p><small className="text-muted">{count} rows</small></p>}
				</div>
			</div>
		)
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
						<th scope="col">Uri</th>
						<th scope="col">Text01</th>
						<th scope="col">Text02</th>
						<th scope="col">Text03</th>
						<th scope="col">Text04</th>
						<th scope="col">
							<img src="/icons/plus-circle.svg" role="img" alt="Plus circle"></img>
						</th>
					</tr>
				</thead>
				<tbody>
					{this.props.items.map(item => <TableRowItemSingleLine key={item.uuid} item={item} />)}
				</tbody>
			</table>
		)
	}
}

class TableRowItemSingleLine extends React.Component {
	render() {
		return (
			<tr>
				<td>{this.props.item.uri}</td>
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
