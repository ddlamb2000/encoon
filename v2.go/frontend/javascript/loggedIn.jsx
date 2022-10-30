// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class LoggedIn extends React.Component {
	constructor(props) {
		super(props)
		const token = localStorage.getItem(`access_token_${dbName}`)
		const payload = this.parseJwt(token)
		this.state = {
			token: token,
			user: payload.user,
			userUuid: payload.userUuid,
			userFirstName: payload.userFirstName,
			userLastName: payload.userLastName,
			error: false,
			disconnect: false,
			isLoaded: false,
			isLoading: false,
			message: "",
			items: [],
		}
	}
  
	parseJwt(token) {
		const base64Url = token.split('.')[1]
		const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
		const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
			return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
		}).join(''))
		const parsedJsonPayload = JSON.parse(jsonPayload)
		return parsedJsonPayload
	}

	componentDidMount() {
		this.setState({isLoading: true})
		const uri = `/${dbName}/api/v1/${gridUri !== "" ? gridUri : 'users'}${uuid !== "" ? '/' + uuid : ''}`
		fetch(uri, {
			headers: {
			'Accept': 'application/json',
			'Content-Type': 'application/json',
			'Authorization': 'Bearer ' + this.state.token
			}
		})
		.then(res => res.json())
		.then(	
			(result) => {
				this.setState({
					isLoading: false,
					isLoaded: true,
					items: result.items,
					error: result.error,
					message: result.message,
					disconnect: result.disconnect
				})
			},
			(error) => {
				this.setState({
					isLoading: false,
					isLoaded: false,
					items: [],
					message: `Something happened: ${error}.`,
					error: true
				})
			}
		)
	}

	render() {
		const { items, isLoading, isLoaded, error, disconnect } = this.state
		if(error) {
			alert(`${dbName}: ${this.state.message}`)
			if(disconnect) {
				localStorage.removeItem(`access_token_${dbName}`)
				location.reload()
			}
			return
		}
		return (
			<div className="container-fluid">
				<Navigation user={this.state.user} userFirstName={this.state.userFirstName} userLastName={this.state.userLastName} />
				<h3>{gridUri}</h3>
				{isLoading && <div>Loading…</div>}
				{isLoaded && items == undefined && <div>No data</div>}
				{isLoaded && items != undefined && uuid != "" && <span className="text-muted">{items[0].uuid}</span>}
				{isLoaded && items != undefined && uuid == "" && <TableRows items={items} />}
				{isLoaded && items != undefined && uuid != "" && <TableSingleRow item={items[0]} />}
			</div>
		)
	}
}

const TableRows = ({ items }) =>
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
		{items.map(item => <TableRowItemSingleLine key={item.uuid} item={item} />)}
	</tbody>
	</table>


const TableRowItemSingleLine = ({ item }) =>
	<tr>
		<td>{item.uri}</td>
		<td>{item.text01}</td>
		<td>{item.text02}</td>
		<td>{item.text03}</td>
		<td>{item.text04}</td>
		<td scope="row">
			<a href={`/${dbName}/${gridUri}/${item.uuid}`}>{item.uuid}</a>
		</td>
	</tr>

const TableSingleRow = ({ item }) =>
	<table className="table table-hover table-sm">
		<thead className="table-light"></thead>
		<tbody>
			<tr><td>Uri</td><td>{item.uri}</td></tr>
			<tr><td>Text01</td><td>{item.text01}</td></tr>
			<tr><td>Text02</td><td>{item.text02}</td></tr>
			<tr><td>Text03</td><td>{item.text03}</td></tr>
			<tr><td>Text04</td><td>{item.text04}</td></tr>
		</tbody>
	</table>
