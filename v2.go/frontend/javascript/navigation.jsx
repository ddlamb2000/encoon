// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class Navigation extends React.Component {
  render() {
    return (
      <nav className="navbar navbar-expand-lg bg-light">
          <div className="container-fluid">
              <a className="navbar-brand">{appName}</a>
              <button className="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarNavAltMarkup" aria-controls="navbarNavAltMarkup" aria-expanded="false" aria-label="Toggle navigation">
                  <span className="navbar-toggler-icon"></span>
              </button>
              <div className="collapse navbar-collapse" id="navbarNavAltMarkup">
                  <div className="navbar-nav">
                      <a className="nav-link" href={`/${dbName}`}>{dbName}</a>
                      <a className="nav-link" href={`/${dbName}/users`}>Users</a> 
                  </div>
              </div>
          </div>
          <div className="navbar-nav">
            <a className="nav-link" href="" onClick={this.logout}>Log&nbsp;out</a>
          </div>
      </nav>
    );
  }

  logout() {
    localStorage.removeItem(`access_token_${dbName}`);
    location.reload();
  }
}