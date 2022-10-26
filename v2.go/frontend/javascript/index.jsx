// εncooη : data structuration, presentation and navigation.
// Copyright David Lambert 2022

class App extends React.Component {
  constructor(props) {
    super(props);
    this.setState();
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
    if (this.loggedIn) return <LoggedIn />;
    return <NotLogged />;
  }
}

const bodyRootElement = document.getElementById("bodyRoot");
const dbName = bodyRootElement.getAttribute("dbName");
const appName = bodyRootElement.getAttribute("appName");
const uuid = bodyRootElement.getAttribute("uuid");
const rootElement = document.getElementById("application");
const root = ReactDOM.createRoot(rootElement);
root.render(<App />);