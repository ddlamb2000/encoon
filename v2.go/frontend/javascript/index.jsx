// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class App extends React.Component {
	constructor(props) {
		super(props)
		this.token = localStorage.getItem(`access_token_${this.props.dbName}`)
		if(this.token) {
			const payload = this.parseJwt(this.token)
			this.user = payload.user
			this.userUuid = payload.userUuid
			this.userFirstName = payload.userFirstName
			this.userLastName = payload.userLastName
			this.loggedIn = true
		}
		else this.loggedIn = false
	}

	render() {
		if(!this.loggedIn) return <Login appName={this.props.appName} dbName={this.props.dbName} />
		if(this.props.gridUri != "") return (
			<div className="container-fluid">
				<Navigation appName={this.props.appName} 
							dbName={this.props.dbName} 
							user={this.user}
							userFirstName={this.userFirstName}
							userLastName={this.userLastName} />
				<Grid token={this.token} dbName={this.props.dbName} gridUri={this.props.gridUri} uuid={this.props.uuid} />
			</div>			
		)
		return (
			<div className="container-fluid">
				<Navigation appName={this.props.appName} 
							dbName={this.props.dbName} 
							user={this.user}
							userFirstName={this.userFirstName}
							userLastName={this.userLastName} />
				<Grid token={this.token} dbName={this.props.dbName} gridUri="users" />
				<Grid token={this.token} dbName={this.props.dbName} gridUri="grids" />
				<Grid token={this.token} dbName={this.props.dbName} gridUri="columns" />
				<Grid token={this.token} dbName={this.props.dbName} />
				<Grid token={this.token} />
				<Grid />
			</div>
		)		
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
}

const rootElement = document.getElementById("application")
const root = ReactDOM.createRoot(rootElement)
root.render(
	<App 
		appName={rootElement.getAttribute("appName")}
		dbName={rootElement.getAttribute("dbName")}
		gridUri={rootElement.getAttribute("gridUri")}
		uuid={rootElement.getAttribute("uuid")}
	/>
)