// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Login extends React.Component {
	constructor(props) {
		super(props)
		this.dbNameInput = null
		this.setDbNameRef = element => { this.dbNameInput = element }
		this.idInput = null
		this.setIdRef = element => { this.idInput = element }
		this.passwordInput = null
		this.setPasswordRef = element => { this.passwordInput = element }

		this.authenticate = () => {
			const updatedDbName = dbName != '' ? dbName : this.dbNameInput.value
			if(updatedDbName == '') {
				alert("Database name is required.")
				return null
			}
			if(this.idInput.value == '' || this.passwordInput.value == '') {
				alert("Username and passphrase are both required.")
				return null
			}
			fetch(`/${updatedDbName}/api/v1/authentication`, {
				method: 'POST',
				headers: { 'Accept': 'application/json', 'Content-Type': 'application/json' },
				body: JSON.stringify({ id: this.idInput.value, password: btoa(this.passwordInput.value) })
			})
			.then( (response) => {
				if(response.status == 400) {
					alert("Incorrect user credentials.")
					return null
				}
				if(response.status != 200) {
					alert(`Problem ${response.status} is reported.`)
					return null
				}
				return response.json() })
				.then( (responseJson) => {
					if (responseJson != null) {
						localStorage.setItem(`access_token_${updatedDbName}`, responseJson.token)
					}
					if(dbName == '') location.href = `/${updatedDbName}/`
					else location.reload()
				} 
			)  
		}
	}

	render() {
		return (
			<div className="container">
				<div className="row">
					<div className="col"></div>
					<div className="col">
						<br/><br/><br/><br/><br/><br/><br/>
						<center>
							<h1>{appName}</h1>
							<small className="text-muted">Data structuration, presentation and navigation.</small>
							<br/>
							<br/>
							{dbName != "" && <h3>{dbName}</h3>}
						</center>
						{dbName == "" &&
							<div className="mb-3">
								<label htmlFor="dbName" className="form-label">Database</label>
								<input type="text" 
									className="form-control" 
									id="dbName" 
									ref={this.setDbNameRef}/>
							</div>
						}
						<div className="mb-3">
							<label htmlFor="id" className="form-label">Username</label>
							<input type="text" 
								className="form-control" 
								id="id" 
								ref={this.setIdRef}
								aria-describedby="idHelp"/>
							<div id="idHelp" className="form-text">We'll never share your data with anyone else.</div>
						</div>
						<div className="mb-3">
							<label htmlFor="password" className="form-label">Passphrase</label>
							<input type="password" 
								className="form-control" 
								id="password"
								ref={this.setPasswordRef}/>
							<div id="passwordHelp" className="form-text">
								A passphrase is a sequence of words or other text used to control access 
								to a computer system, program or data. 
								It is similar to a password in usage, but a passphrase is generally longer 
								for added security.
							</div>
						</div>
						<div className="d-grid gap-2">
							<button type="button" 
								className="btn btn-outline-primary"
								onClick={this.authenticate}>
								Log in {dbName} <img src="/icons/box-arrow-in-right.svg" role="img" alt="Log in" />
							</button>
						</div>
					</div>
					<div className="col"></div>
				</div>
			</div>
		)
	}
}