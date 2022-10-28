// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class NotLogged extends React.Component {
    constructor(props) {
      super(props);
      this.authenticate = this.authenticate.bind(this);
    }
  
    authenticate() {
      let id = document.getElementById("id").value;
      let password = document.getElementById("password").value;

      fetch(`/${dbName}/api/v1/authentication`, {
        method: 'POST',
        headers: {
          'Accept': 'application/json',
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({id: id, password: btoa(password)})
      })
      .then( (response) => {
        if (response.status == 400) {
          alert("Invalid ID or password.")
          return null;
        }
        if (response.status != 200) {
          alert('Looks like there was a problem. Status Code: ' + response.status);
          return null;
        }
        return response.json() })
      .then( (responseJson) => {
        if (responseJson != null) {
          localStorage.setItem(`access_token_${dbName}`, responseJson.token);
        }
        location.reload();
      } )
    }
    
    render() {
      return (
        <div className="container">
          <div className="row">
            <div className="col"></div>
            <div className="col">
              <br/>
              <br/>
              <br/>
              <br/>
              <center><h1>{appName}</h1></center>
              <center><h4>{dbName}</h4></center>
              <br/>
              <div className="mb-3">
                <label htmlFor="id" className="form-label">Identifier</label>
                <input type="email" className="form-control" id="id" aria-describedby="emailHelp"/>
                <div id="emailHelp" className="form-text">We'll never share your data with anyone else.</div>
              </div>
              <div className="mb-3">
                <label htmlFor="password" className="form-label">Password</label>
                <input type="password" className="form-control" id="password"/>
              </div>
              <button type="submit" className="btn btn-primary" onClick={this.authenticate}>Submit</button>
            </div>
            <div className="col"></div>
          </div>
        </div>
      );
    }
  }