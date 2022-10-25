class App extends React.Component {
  constructor(props) {
    super(props);
    this.setup();
    this.setState();
  }

  setup() {
    $.ajaxSetup({
      beforeSend: (r) => {
        if (localStorage.getItem("access_token")) {
          r.setRequestHeader(
            "Authorization",
            "Bearer " + localStorage.getItem("access_token")
          );
        }
      }
    });
  }

  setState() {
    let idToken = localStorage.getItem("access_token");
    if (idToken) {
      this.loggedIn = true;
    } else {
      this.loggedIn = false;
    }
  }

  render() {
    if (this.loggedIn) {
      return <LoggedIn />;
    }
    return <NotLogged />;
  }
}

class NotLogged extends React.Component {
  constructor(props) {
    super(props);
    this.authenticate = this.authenticate.bind(this);
  }

  authenticate() {
    // this.WebAuth = new auth0.WebAuth({
    //   domain: AUTH0_DOMAIN,
    //   clientID: AUTH0_CLIENT_ID,
    //   scope: "openid profile",
    //   audience: AUTH0_API_AUDIENCE,
    //   responseType: "token id_token",
    //   redirectUri: AUTH0_CALLBACK_URL
    // });
    // this.WebAuth.authorize();

    this.serverRequest();
  }

  serverRequest() {
    var id = document.getElementById("id").value;
    var password = document.getElementById("password").value;

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
      if (response.status != 200) {
        console.log('Looks like there was a problem. Status Code: ' +
          response.status);
        alert(response.status)
        return null;
      }
      return response.json() })
    .then( (responseJson) => {
      if (responseJson != null) {
        localStorage.setItem("access_token", responseJson.token);
      }
      location.reload();
      return;
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

class LoggedIn extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
        error: null,
        isLoaded: false,
        items: [],
    };
  }

  render() {
    const { items, isLoaded, error } = this.state;

    if (error) {
        return <div>Erreur : {error.message}</div>;
    } else if (!isLoaded) {
        return <div>Chargement…</div>;
    } else {
        if(uuid !== "") {
            return (
                <div>
                    <h2>User</h2>
                    <table className="table table-hover table-sm">
                        <thead className="table-light">
                        </thead>
                        <tbody>
                            <tr><td>Uuid</td><td>{items.uuid}</td></tr>
                            <tr><td>Email</td><td>{items.email}</td></tr>
                            <tr><td>First Name</td><td>{items.firstName}</td></tr>
                            <tr><td>Last Name</td><td>{items.lastName}</td></tr>
                        </tbody>
                    </table>
                    <span className="pull-right">
                      <a onClick={this.logout}>Log out</a>
                    </span>
                </div>
            );
        }
        else {
            return (
                <div>
                    <h2>Users</h2>
                    <table className="table table-hover table-sm">
                        <thead className="table-light">
                            <tr>
                                <th scope="col">Email</th>
                                <th scope="col">First Name</th>
                                <th scope="col">Last Name</th>
                                <th scope="col">
                                    <svg xmlns="http://www.w3.org/2000/svg" width="16" height="16" fill="currentColor" className="bi bi-battery-full" viewBox="0 0 16 16">
                                        <path d="M2 6h10v4H2V6z"/>
                                        <path d="M2 4a2 2 0 0 0-2 2v4a2 2 0 0 0 2 2h10a2 2 0 0 0 2-2V6a2 2 0 0 0-2-2H2zm10 1a1 1 0 0 1 1 1v4a1 1 0 0 1-1 1H2a1 1 0 0 1-1-1V6a1 1 0 0 1 1-1h10zm4 3a1.5 1.5 0 0 1-1.5 1.5v-3A1.5 1.5 0 0 1 16 8z"/>
                                    </svg>
                                </th>
                            </tr>
                        </thead>
                        <tbody>
                            {items.map(item => (
                                <tr key={item.uuid}>
                                    <td>{item.email}</td>
                                    <td>{item.firstName}</td>
                                    <td>{item.lastName}</td>
                                    <th scope="row"><a href={`/${dbName}/user/${item.uuid}`}>{item.uuid}</a></th>
                                </tr>
                            ))}
                        </tbody>
                    </table>
                    <span className="pull-right">
                      <a onClick={this.logout}>Log out</a>
                    </span>
                </div>
            );
        }
    }
  }

  logout() {
    localStorage.removeItem("access_token");
    location.reload();
  }

  componentDidMount() {
    if(uuid !== "") {
        fetch(`/${dbName}/api/v1/user/${uuid}`, {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem("access_token")
          }
        })
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    isLoaded: true,
                    items: result.user
                });
            },
            (error) => {
                this.setState({
                    isLoaded: true,
                    error
                });
            }
        )
    }
    else {
        fetch(`/${dbName}/api/v1/users`, {
          headers: {
            'Accept': 'application/json',
            'Content-Type': 'application/json',
            'Authorization': 'Bearer ' + localStorage.getItem("access_token")
          }
        })
        .then(res => res.json())
        .then(
            (result) => {
                this.setState({
                    isLoaded: true,
                    items: result.users
                });
            },
            (error) => {
                this.setState({
                    isLoaded: true,
                    error
                });
            }
        )
    }
  }
}

const bodyRootElement = document.getElementById("bodyRoot");
const dbName = bodyRootElement.getAttribute("dbName");
const appName = bodyRootElement.getAttribute("appName");
const uuid = bodyRootElement.getAttribute("uuid");
const rootElement = document.getElementById("application");
const root = ReactDOM.createRoot(rootElement);
root.render(<App />);
