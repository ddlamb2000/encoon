// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Grid extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			error: false,
			disconnect: false,
			isLoaded: false,
			isLoading: false,
			items: [],
			count: 0,
		}
	}
  
	componentDidMount() {
		if(this.props.gridUri == undefined) {
			this.setState({error: "Missing parameter gridUri"})
			return
		}
		if(this.props.uuid == undefined) {
			this.setState({error: "Missing parameter uuid"})
			return
		}

		this.setState({isLoading: true})
		const uri = `/${this.props.dbName}/api/v1/${this.props.gridUri}${this.props.uuid != "" ? '/' + this.props.uuid : ''}`
		console.log(`Trigger ${uri}`)
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
					(result) => this.setState({
							isLoading: false,
							isLoaded: true,
							items: result.items,
							count: result.count,
							error: result.error,
							disconnect: result.disconnect
					}),
					(error) => this.setState({
							isLoading: false,
							isLoaded: false,
							items: [],
							error: error.message
					})
				)
			} else {
				alert("ooops")
			}
		})
	}

	render() {
		const { isLoading, isLoaded, error, disconnect, items, count } = this.state
		if(error && disconnect) {
			alert(error)
			localStorage.removeItem(`access_token_${this.props.dbName}`)
			location.reload()
			return
		}
		return (
			<div>
				<h3>{this.props.gridUri}{isLoading && <Spinner />}</h3>
				{error && !isLoading && !isLoaded && <div className="alert alert-primary" role="alert">{error}</div>}
				{isLoaded && items && count == 0 && <div className="alert alert-secondary" role="alert">No data</div>}
				{isLoaded && items && count > 0 && this.props.uuid == "" && <TableRows items={items} />}
				{isLoaded && items && count > 0 && this.props.uuid != "" && <TableSingleRow item={items[0]} />}
			</div>
		)
	}
}

// Grid.defaultProps = {
// 	uuid: ''
// }

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
				<span className="text-muted">{this.props.item.uuid}</span>
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
