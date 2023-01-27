// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2023

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
