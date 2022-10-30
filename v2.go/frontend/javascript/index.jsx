// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class App extends React.Component {
	constructor(props) {
		super(props)
		this.token = localStorage.getItem(`access_token_${dbName}`)
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
		if (!this.loggedIn) return <Login />
		return (
			<div className="container-fluid">
				<Navigation user={this.user} userFirstName={this.userFirstName} userLastName={this.userLastName} />
				<Grid token={this.token} />
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

const bodyRootElement = document.getElementById("bodyRoot")
const dbName = bodyRootElement.getAttribute("dbName")
const appName = bodyRootElement.getAttribute("appName")
const uuid = bodyRootElement.getAttribute("uuid")
const gridUri = bodyRootElement.getAttribute("gridUri")
const rootElement = document.getElementById("application")
const root = ReactDOM.createRoot(rootElement)
root.render(<App />)