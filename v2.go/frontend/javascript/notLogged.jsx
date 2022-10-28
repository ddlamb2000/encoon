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
        body: JSON.stringify({
          id: id,
          password: password,
        })
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
            <div className="col-xs-4 col-xs-offset-4 jumbotron text-center">
              <h1>{appName}</h1>
              <h3>{dbName}</h3>
              <br />
              <p>Sign in to get access </p>
              <br />
              <div className="form-group has-feedback">
                <input type="text" name="id" id="id" size="36" placeholder="ID"/>
                <span className="glyphicon glyphicon-envelope form-control-feedback"></span>
              </div>
              <br />
              <div className="form-group has-feedback">
                <input type="password" name="password" id="password" size="36" placeholder="Password"/>
                <span className="glyphicon glyphicon-lock form-control-feedback"></span>
              </div>
              <br />
              <a onClick={this.authenticate}
                 className="btn btn-primary btn-lg btn-login btn-block">
                Sign In
              </a>
            </div>
          </div>
        </div>
      );
    }
  }