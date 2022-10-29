// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class NotLogged extends React.Component {
  constructor(props) {
    super(props)
    this.idInput = null
    this.setIdRef = element => { this.idInput = element }
    this.passwordInput = null
    this.setPasswordRef = element => { this.passwordInput = element }

    this.authenticate = () => {
      fetch(`/${dbName}/api/v1/authentication`, {
        method: 'POST',
        headers: { 'Accept': 'application/json', 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: this.idInput.value, password: btoa(this.passwordInput.value) })
      })
      .then( (response) => {
        if(response.status == 400) {
          alert("Incorrect credentials.")
          return null
        }
        if(response.status != 200) {
          alert(`Problem ${response.status} is reported.`)
          return null
        }
        return response.json() })
      .then( (responseJson) => {
        if (responseJson != null) {
          localStorage.setItem(`access_token_${dbName}`, responseJson.token)
        }
        location.reload()
      } )  
    }
  }

  render() {
    return (
      <div className="container">
        <div className="row">
          <div className="col"></div>
          <div className="col">
            <br/><br/><br/>
            <center><h1>{appName}</h1></center>
            <center><h4>{dbName}</h4></center>
            <br/>
            <div className="mb-3">
              <label htmlFor="id" className="form-label">Identifier</label>
              <input
                type="text" 
                className="form-control" 
                id="id" 
                ref={this.setIdRef}
                aria-describedby="idHelp"/>
              <div id="idHelp" className="form-text">We'll never share your data with anyone else.</div>
            </div>
            <div className="mb-3">
              <label htmlFor="password" className="form-label">Password</label>
              <input 
                type="password" 
                className="form-control" 
                id="password"
                ref={this.setPasswordRef}/>
            </div>
            <button
              type="button" 
              className="btn btn-primary" 
              onClick={this.authenticate}>Log in
            </button>
          </div>
          <div className="col"></div>
        </div>
      </div>
    );
  }
}