// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

class App extends React.Component {
	constructor(props) {
		super(props)
		this.token = localStorage.getItem(`access_token_${this.props.dbName}`)
		if(this.token != "") this.verifyToken()
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
								<div>
									<Grid token={this.token} dbName={this.props.dbName} gridUuid={UuidUsers} navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
									<Grid token={this.token} dbName={this.props.dbName} gridUuid={UuidGrids} navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
									<Grid token={this.token} dbName={this.props.dbName} gridUuid={UuidColumns} navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
									<Grid token={this.token} dbName={this.props.dbName} navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
									<Grid token={this.token} navigateToGrid={(gridUuid, uuid) => this.navigateToGrid(gridUuid, uuid)} />
									<Grid />
								</div>
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

class Login extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			isLoading: false,
			errorDb: "",
			errorId: "",
			errorPassword: "",
			error: ""
		}
		this.dbNameInput = null
		this.setDbNameRef = element => { this.dbNameInput = element }
		this.idInput = null
		this.setIdRef = element => { this.idInput = element }
		this.passwordInput = null
		this.setPasswordRef = element => { this.passwordInput = element }
	}

	authenticate() {
		this.setState({
			isLoading: true,
			errorDb: "",
			errorId: "",
			errorPassword: "",
			error: ""
		})
		const updatedDbName = this.getUpdatedDbName()
		if(updatedDbName == '') this.setState({ errorDb: "Please enter a database name." })
		if(this.idInput.value == '') this.setState({ errorId: "Please enter a username." })
		if(this.passwordInput.value == '') this.setState({ errorPassword: "Please enter a passphrase." })
		if(updatedDbName == '' || this.idInput.value == '' || this.passwordInput.value == '') {
			this.setState({ isLoading: false })
			return false
		}
		fetch(`/${updatedDbName}/api/v1/authentication`, {
			method: 'POST',
			headers: { 'Accept': 'application/json', 'Content-Type': 'application/json' },
			body: JSON.stringify({ id: this.idInput.value, password: btoa(this.passwordInput.value)})
		})
		.then(response => this.handleAuthenticateResponse(response))
	}

	getUpdatedDbName() {
		return this.props.dbName != '' ? this.props.dbName : this.dbNameInput.value
	}

	handleAuthenticateResponse(response) {
		const updatedDbName = this.getUpdatedDbName()
		this.setState({ isLoading: false })
		const contentType = response.headers.get("content-type")
		if(contentType && contentType.indexOf("application/json") != -1) {
			return response.json().then(	
				(result) => {
					if(response.status == 200) {
						localStorage.setItem(`access_token_${updatedDbName}`, result.token)
						if(this.props.dbName == '') location.href = `/${updatedDbName}/`
						else location.reload()
						return true
					}
					if(response.status == 401) {
						this.setState({ errorPassword: "Invalid username or passphrase" })
						return false
					}
					this.setState({ error: result.error })
				},
				(error) => { this.setState({ error: error.message }) }
			)
		} else { this.setState({ error: `[${response.status}] Internal server issue.` }) }
		return false

	}

	render() {
		const { isLoading, errorDb, errorId, errorPassword, error } = this.state
		const variantDb = errorDb ? "form-control is-invalid" : "form-control"
		const variantId = errorId ? "form-control is-invalid" : "form-control"
		const variantPassword = errorPassword ? "form-control is-invalid" : "form-control"
		if(trace) console.log("[Login.render()]")
		return (
			<div className="container">
				<div className="row">
					<div className="col"></div>
					<form className="col">
						<br/><br/><br/><br/><br/><br/><br/>
						<center>
							<h1>{this.props.appName}</h1>
							<small className="text-muted">{this.props.appTag}</small>
							<br/><br/>
							{this.props.dbName != "" && <h3>{this.props.dbName}</h3>}
						</center>
						{this.props.dbName == "" &&
							<div className="mb-3">
								<label htmlFor="dbName" className="form-label">Database</label>
								<input type="text" 
									className={variantDb}
									id="dbName" 
									ref={this.setDbNameRef} />
								<div className="invalid-feedback">{errorDb}</div>
							</div>
						}
						<div className="mb-3">
							<label htmlFor="id" className="form-label">Username</label>
							<input type="text" 
								className={variantId}
								id="id" 
								ref={this.setIdRef}
								aria-describedby="idHelp" />
							<div className="invalid-feedback">{errorId}</div>
							<div id="idHelp" className="form-text">We'll never share your data with anyone else.</div>
						</div>
						<div className="mb-3">
							<label htmlFor="password" className="form-label">Passphrase</label>
							<input type="password" 
								className={variantPassword}
								id="password"
								ref={this.setPasswordRef} />
							 <div className="invalid-feedback">{errorPassword}</div>
							<div id="passwordHelp" className="form-text">
								A passphrase is a sequence of words or other text used to control access 
								to a computer system, program or data. 
								It is similar to a password in usage, but a passphrase is generally longer 
								for added security.
							</div>
						</div>
						{error && <div className="mb-3 alert alert-danger" role="alert">{error}</div>}
						<div className="d-grid gap-2">
							<button type="button" className="btn btn-outline-primary" onClick={() => this.authenticate()}>
								Log in {this.props.dbName}&nbsp;
								{!isLoading && <i className="bi bi-box-arrow-in-right"></i>}
								{isLoading && <Spinner />}
							</button>
						</div>
					</form>
					<div className="col"></div>
				</div>
			</div>
		)
	}
}

class Header extends React.Component {
	render() {
		if(trace) console.log("[Header.render()]")
		return (
            <header className="navbar sticky-top bg-light flex-md-nowrap p-0 shadow">
                <a className="navbar-brand col-md-3 col-lg-2 me-0 px-3 fs-6" href={"/" + this.props.dbName}>{this.props.appName} | {this.props.dbName}</a>
                <div className="navbar-text">
                    <small className="nav-item text-nowrap px-4 text-muted">{this.props.appTag}</small>
                </div>
                <input className="form-control form-control-dark w-100 rounded-0 border-0" type="text" placeholder="Search" aria-label="Search" />
                <div className="navbar-text">
                    <div className="nav-item text-nowrap px-4">{this.props.user}</div>
                </div>
                <div className="navbar-nav">
                    <div className="nav-item text-nowrap">
                        <a className="nav-link px-3">
                            <button type="button"
                                    className="btn btn-outline-secondary btn-sm"
                                    onClick={
                                        () => {
                                            localStorage.removeItem(`access_token_${this.props.dbName}`)
                                            location.reload()
                                        }
                                    }>
                                Log out <i className="bi bi-box-arrow-right"></i>
                            </button>
                        </a>
                    </div>
                </div>
            </header>
		)
	}
}

class DateTime extends React.Component {
	constructor(props) {
		super(props)
		this.state = {
			timeAgo: this.getTimeAgo()
		}
	}

	getTimeAgo() {
		const localNow = new Date
		const localDate = new Date(this.props.dateTime)
		const localNowUTC =  Date.UTC(localNow.getUTCFullYear(),
										localNow.getUTCMonth(),
										localNow.getUTCDate(),
										localNow.getUTCHours(),
										localNow.getUTCMinutes(),
										localNow.getUTCSeconds())
		const localDateUTC =  Date.UTC(localDate.getUTCFullYear(),
										localDate.getUTCMonth(),
										localDate.getUTCDate(),
										localDate.getUTCHours(),
										localDate.getUTCMinutes(),
										localDate.getUTCSeconds())
		const seconds = (localNowUTC - localDateUTC) / 1000
		const MINUTE = 60, HOUR = MINUTE * 60, DAY = HOUR * 24, WEEK = DAY * 7, MONTH = DAY * 30, YEAR = DAY * 365		
		if(seconds < MINUTE) return `${Math.round(seconds)} sec ago`
		if(seconds < HOUR) return `${Math.round(seconds / MINUTE)} min ago`
		if(seconds < DAY) return `${Math.round(seconds / HOUR)} hour ago`
		if(seconds < WEEK) return `${Math.round(seconds / DAY)} day ago`
		if(seconds < MONTH) return `${Math.round(seconds / WEEK)} week ago`
		if(seconds < YEAR) return `${Math.round(seconds / MONTH)} month ago`
		return `${Math.round(seconds / YEAR)} year ago`
	}

	componentDidMount() {
		this.timerID = setInterval(() => { this.setState({ timeAgo: this.getTimeAgo() }) }, 1000)
	}

	componentWillUnmount() {
		clearInterval(this.timerID)
	}

	render() {
		return (
			<span>{this.props.dateTime} <small><em>({this.state.timeAgo})</em></small></span>
		)
	}
}

class Spinner extends React.Component {
	render() {
		return (
			<span className="spinner-grow spinner-grow-sm ms-1" role="status">
				<span className="visually-hidden">Loading...</span>
			</span>
		)
	}
}

class Navigation extends React.Component {
	render() {
		if(trace) console.log("[Navigation.render()]")
		return (
			<nav className="position-sticky pt-4 sidebar-sticky">
				<Grid token={this.props.token}
						dbName={this.props.dbName}
						gridUuid={UuidGrids}
						navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
						innerGrid={true}
						miniGrid={true}
						gridTitle="My grids"
						gridSubTitle="Grids I own"
						filterColumnOwned={true}
						filterColumnName='relationship3'
						filterColumnGridUuid={UuidGrids}
						filterColumnValue={this.props.userUuid} />
				<Grid token={this.props.token}
						dbName={this.props.dbName}
						gridUuid={UuidGrids}
						navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
						innerGrid={true}
						miniGrid={true}
						gridTitle="Edit grids"
						gridSubTitle="Grids I can edit"
						filterColumnOwned={true}
						filterColumnName='relationship5'
						filterColumnGridUuid={UuidGrids}
						filterColumnValue={this.props.userUuid}
						noEdit={true} />
				<Grid token={this.props.token}
						dbName={this.props.dbName}
						gridUuid={UuidGrids}
						navigateToGrid={(gridUuid, uuid) => this.props.navigateToGrid(gridUuid, uuid)}
						innerGrid={true}
						miniGrid={true}
						gridTitle="View grids"
						gridSubTitle="Grids I can view"
						filterColumnOwned={true}
						filterColumnName='relationship4'
						filterColumnGridUuid={UuidGrids}
						filterColumnValue={this.props.userUuid}
						noEdit={true} />
			</nav>
		)
	}
}

const rootElement = document.getElementById("application")
const root = ReactDOM.createRoot(rootElement)

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

root.render(
	<App  appName={rootElement.getAttribute("appName")}
			appTag={rootElement.getAttribute("appTag")}
			dbName={rootElement.getAttribute("dbName")}
			gridUuid={rootElement.getAttribute("gridUuid")}
			uuid={rootElement.getAttribute("uuid")} />
)