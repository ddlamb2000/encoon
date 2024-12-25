// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class App extends React.Component {
	constructor(props) {
		super(props)
		this.token = localStorage.getItem(`access_token_${this.props.dbName}`)
		if(this.token != null && this.token != "") this.verifyToken()
		else this.loggedIn = false
		this.state = { gridUuid: this.props.gridUuid, uuid: this.props.uuid }
	}

	verifyToken() {
		const payload = this.parseJwt(this.token)
		this.user = payload.user
		this.userUuid = payload.userUuid
		this.userFirstName = payload.userFirstName
		this.userLastName = payload.userLastName
		if(payload.expires == "") {
			if(trace) console.log("No token expiration date.")
			this.loggedIn = false
		}
		else this.verifyTokenExpiration(payload)
	}

	verifyTokenExpiration(payload) {
		const expires = new Date(Date.parse(payload.expires))
		if (expires == "Invalid Date") {
			if(trace) console.log("Invalid token expiration date.")
			this.loggedIn = false
		}
		else {
			const now = new Date()
			if(now > expires) {
				if(trace) console.log("Token expired.")
				this.loggedIn = false
			}
			else if(now < expires) this.loggedIn = true
		}
	}

	render() {
		if(!this.loggedIn) return <Login appName={this.props.appName} appTag={this.props.appTag} dbName={this.props.dbName} />
		const gridUuid = this.state.gridUuid
		const uuid = this.state.uuid
		if(trace) console.log("[App.render()] gridUuid=", gridUuid, ", uuid=", uuid)
		return (
			<div>
				<Header appName={this.props.appName} 
						appTag={this.props.appTag}
						dbName={this.props.dbName} 
						user={this.user}
						userFirstName={this.userFirstName}
						userLastName={this.userLastName} />
				<div className="container-fluid">
					<div className="row">
						<nav id="sidebarMenu" className="col-md-3 col-lg-2 d-md-block bg-light sidebar collapse">
							<Navigation appName={this.props.appName} 
										appTag={this.props.appTag}
										dbName={this.props.dbName} 
										user={this.user}
										userUuid={this.userUuid}
										userFirstName={this.userFirstName}
										userLastName={this.userLastName}
										token={this.token}
										navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
						</nav>
						<main className="col-md-9 ms-sm-auto col-lg-10 px-md-4">
							{gridUuid == "" &&
								<Grid token={this.token}
										dbName={this.props.dbName}
										gridUuid={UuidGrids}
										navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
							}
							{gridUuid != "" &&
								<Grid token={this.token}
										dbName={this.props.dbName}
										gridUuid={gridUuid}
										uuid={uuid}
										navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
							}
						</main>
					</div>
				</div>
			</div>
		)
	}

	parseJwt(token) {
		try {
			const base64Url = token.split('.')[1]
			const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/')
			const jsonPayload = decodeURIComponent(window.atob(base64).split('').map(function(c) {
				return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2)
			}).join(''))
			const parsedJsonPayload = JSON.parse(jsonPayload)
			return parsedJsonPayload
		} catch (error) {
			console.error("Invalid token", error)
		}
		return ""
	}

	navigateToGrid(gridUuid, uuid) {
		if(trace) console.log("[App.navigateToGrid()] gridUuid=", gridUuid, ", uuid=", uuid)
		const url = `/${this.props.dbName}/${this.state.gridUuid}` + (this.state.uuid == "" ? "" : `/${this.state.uuid}`)
		history.replaceState({ gridUuid: this.state.gridUuid, uuid: this.state.uuid }, null, url)
		this.setState({ gridUuid: gridUuid, uuid: uuid })
	}

	componentDidMount() {
		window.addEventListener('popstate', (e) => {
			e.preventDefault()
			if(e && e.isTrusted && e.state != null) {
				if(trace) console.log("[App.componentDidMount()] popstate, e=", e)
				this.setState({ gridUuid: e.state.gridUuid, uuid: e.state.uuid })
			}
		})
	}
}

const UuidTextColumnType                 = "65f3c258-fb1e-4f8b-96ca-f790e70d29c1"
const UuidIntColumnType                  = "8c28d527-66f4-481c-902e-ac1e65a8abb0"
const UuidReferenceColumnType            = "c8b16312-d4f0-40a5-aa04-c0bc1350fea7"
const UuidPasswordColumnType             = "5f038b21-d9a4-45fc-aa3f-fc405342c287"
const UuidBooleanColumnType              = "6e205ebd-6567-44dc-8fd4-ef6ad281ab40"
const UuidUuidColumnType                 = "d7c004ff-da5e-4a18-9520-cd42b2847508"
const UuidRichTextColumnType   			 = "28ac131f-f04b-4350-b464-3db4f8920597"
const UuidGrids                          = "f35ef7de-66e7-4e51-9a09-6ff8667da8f7"
const UuidGridColumnName                 = "e9e4a415-c31e-4383-ae70-18949d6ec692"
const UuidUsers                          = "018803e1-b4bf-42fa-b58f-ac5faaeeb0c2"
const UuidColumns                        = "533b6862-add3-4fef-8f93-20a17aaaaf5a"

const trace = false

const rootElement = document.getElementById("application")
const root = ReactDOM.createRoot(rootElement)
root.render(
	<App  appName={rootElement.getAttribute("appName")}
			appTag={rootElement.getAttribute("appTag")}
			dbName={rootElement.getAttribute("dbName")}
			gridUuid={rootElement.getAttribute("gridUuid")}
			uuid={rootElement.getAttribute("uuid")} />
)